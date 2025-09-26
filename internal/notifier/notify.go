package notifier

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/wcy-dt/ponghub/internal/types/structures/checker"
	"github.com/wcy-dt/ponghub/internal/types/types/chk_result"
	"github.com/wcy-dt/ponghub/internal/types/types/default_config"
)

// WriteNotifications sends notifications based on the service check results
func WriteNotifications(checkResult []checker.Service, certNotifyDays int) {
	statusNoneEndpoints := collectUnavailableEndpoints(checkResult)
	certProblemEndpoints := collectCertProblemEndpoints(checkResult, certNotifyDays)

	if len(statusNoneEndpoints) == 0 && len(certProblemEndpoints) == 0 {
		// if no endpoints have issues, do nothing
		return
	}

	notifyPath := default_config.GetNotifyPath()
	if err := removeExistingNotifyFile(notifyPath); err != nil {
		return
	}

	f, err := os.Create(notifyPath)
	if err != nil {
		log.Println("Error creating notify file:", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Error closing notify file:", err)
		}
	}()

	writeNotificationReport(f, statusNoneEndpoints, certProblemEndpoints)
}

// collectUnavailableEndpoints finds all endpoints with status NONE
func collectUnavailableEndpoints(checkResult []checker.Service) map[string][]checker.Endpoint {
	statusNoneEndpoints := make(map[string][]checker.Endpoint)
	for _, serviceResult := range checkResult {
		for _, endpointResult := range serviceResult.Endpoints {
			if endpointResult.Status == chk_result.NONE {
				statusNoneEndpoints[serviceResult.Name] = append(statusNoneEndpoints[serviceResult.Name], endpointResult)
			}
		}
	}
	return statusNoneEndpoints
}

// collectCertProblemEndpoints finds all endpoints whose certificates are expired or expiring soon
func collectCertProblemEndpoints(checkResult []checker.Service, certNotifyDays int) map[string][]checker.Endpoint {
	certProblemEndpoints := make(map[string][]checker.Endpoint)
	for _, serviceResult := range checkResult {
		for _, endpointResult := range serviceResult.Endpoints {
			if endpointResult.IsHTTPS && (endpointResult.IsCertExpired || endpointResult.CertRemainingDays <= certNotifyDays) {
				certProblemEndpoints[serviceResult.Name] = append(certProblemEndpoints[serviceResult.Name], endpointResult)
			}
		}
	}
	return certProblemEndpoints
}

// removeExistingNotifyFile removes the existing notify file if it exists
func removeExistingNotifyFile(notifyPath string) error {
	if err := os.Remove(notifyPath); err != nil && !os.IsNotExist(err) {
		log.Println("Error removing notify file:", err)
		return err
	}
	return nil
}

// writeNotificationReport writes the complete notification report to the file
func writeNotificationReport(f *os.File, statusNoneEndpoints, certProblemEndpoints map[string][]checker.Endpoint) {
	writeHeader(f)
	writeUnavailableServices(f, statusNoneEndpoints)
	writeCertificateIssues(f, certProblemEndpoints)
	writeSummary(f, statusNoneEndpoints, certProblemEndpoints)
}

// writeHeader writes the report header with timestamp
func writeHeader(f *os.File) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	writeToFile(f, fmt.Sprintf("=== PongHub Service Status Report ===\n"))
	writeToFile(f, fmt.Sprintf("Generated at: %s\n\n", currentTime))
}

// writeUnavailableServices writes information about unavailable services
func writeUnavailableServices(f *os.File, statusNoneEndpoints map[string][]checker.Endpoint) {
	if len(statusNoneEndpoints) == 0 {
		return
	}

	writeToFile(f, "🔴 UNAVAILABLE SERVICES:\n")
	writeToFile(f, strings.Repeat("=", 50)+"\n")

	for serviceName, endpoints := range statusNoneEndpoints {
		writeToFile(f, fmt.Sprintf("\n📋 Service: %s\n", serviceName))
		for _, endpoint := range endpoints {
			writeUnavailableEndpointDetails(f, endpoint)
		}
	}
}

// writeUnavailableEndpointDetails writes detailed information about an unavailable endpoint
func writeUnavailableEndpointDetails(f *os.File, endpoint checker.Endpoint) {
	writeToFile(f, fmt.Sprintf("  • URL: %s\n", endpoint.URL))
	writeToFile(f, fmt.Sprintf("    Method: %s\n", endpoint.Method))

	if endpoint.StatusCode > 0 {
		writeToFile(f, fmt.Sprintf("    Status Code: %d\n", endpoint.StatusCode))
	}
	if endpoint.ResponseTime > 0 {
		writeToFile(f, fmt.Sprintf("    Response Time: %v\n", endpoint.ResponseTime))
	}

	writeToFile(f, fmt.Sprintf("    Attempts: %d/%d successful\n", endpoint.SuccessNum, endpoint.AttemptNum))
	writeToFile(f, fmt.Sprintf("    Check Time: %s - %s\n", endpoint.StartTime, endpoint.EndTime))

	writeFailureDetails(f, endpoint.FailureDetails)
	writeResponseBody(f, endpoint.ResponseBody)
	writeToFile(f, "\n")
}

// writeFailureDetails writes failure details if available
func writeFailureDetails(f *os.File, failureDetails []string) {
	if len(failureDetails) == 0 {
		return
	}

	writeToFile(f, "    Failure Details:\n")
	for _, detail := range failureDetails {
		writeToFile(f, fmt.Sprintf("      - %s\n", detail))
	}
}

// writeResponseBody writes response body if available and not too long
func writeResponseBody(f *os.File, responseBody string) {
	if len(responseBody) > 0 && len(responseBody) < 500 {
		writeToFile(f, fmt.Sprintf("    Response Body: %s\n", strings.TrimSpace(responseBody)))
	}
}

// writeCertificateIssues writes information about certificate issues
func writeCertificateIssues(f *os.File, certProblemEndpoints map[string][]checker.Endpoint) {
	if len(certProblemEndpoints) == 0 {
		return
	}

	writeToFile(f, "\n🔐 CERTIFICATE ISSUES:\n")
	writeToFile(f, strings.Repeat("=", 50)+"\n")

	for serviceName, endpoints := range certProblemEndpoints {
		writeToFile(f, fmt.Sprintf("\n📋 Service: %s\n", serviceName))
		for _, endpoint := range endpoints {
			writeCertEndpointDetails(f, endpoint)
		}
	}
}

// writeCertEndpointDetails writes detailed information about certificate issues
func writeCertEndpointDetails(f *os.File, endpoint checker.Endpoint) {
	writeToFile(f, fmt.Sprintf("  • URL: %s\n", endpoint.URL))

	writeCertificateStatus(f, endpoint)

	writeToFile(f, fmt.Sprintf("    Days Remaining: %d\n", endpoint.CertRemainingDays))
	if endpoint.StatusCode > 0 {
		writeToFile(f, fmt.Sprintf("    Status Code: %d\n", endpoint.StatusCode))
	}
	if endpoint.ResponseTime > 0 {
		writeToFile(f, fmt.Sprintf("    Response Time: %v\n", endpoint.ResponseTime))
	}
	writeToFile(f, fmt.Sprintf("    Check Time: %s - %s\n", endpoint.StartTime, endpoint.EndTime))
	writeToFile(f, "\n")
}

// writeCertificateStatus writes the certificate status with appropriate emoji and message
func writeCertificateStatus(f *os.File, endpoint checker.Endpoint) {
	if endpoint.IsCertExpired {
		writeToFile(f, "    ❌ Certificate Status: EXPIRED\n")
	} else {
		certStatus := "⚠️  Certificate Status: EXPIRES SOON"
		if endpoint.CertRemainingDays <= 1 {
			certStatus = "🚨 Certificate Status: EXPIRES IN 1 DAY OR LESS"
		}
		writeToFile(f, fmt.Sprintf("    %s\n", certStatus))
	}
}

// writeSummary writes the summary statistics
func writeSummary(f *os.File, statusNoneEndpoints, certProblemEndpoints map[string][]checker.Endpoint) {
	writeToFile(f, "\n📊 SUMMARY:\n")
	writeToFile(f, strings.Repeat("=", 50)+"\n")

	unavailableCount := countEndpoints(statusNoneEndpoints)
	certIssueCount := countEndpoints(certProblemEndpoints)

	writeToFile(f, fmt.Sprintf("Unavailable Endpoints: %d\n", unavailableCount))
	writeToFile(f, fmt.Sprintf("Certificate Issues: %d\n", certIssueCount))
	writeToFile(f, fmt.Sprintf("Total Issues: %d\n", unavailableCount+certIssueCount))
}

// countEndpoints counts the total number of endpoints in the map
func countEndpoints(endpointsMap map[string][]checker.Endpoint) int {
	count := 0
	for _, endpoints := range endpointsMap {
		count += len(endpoints)
	}
	return count
}

// writeToFile is a helper function that writes to file and handles errors
func writeToFile(f *os.File, content string) {
	if _, err := f.WriteString(content); err != nil {
		log.Println("Error writing to notify file:", err)
	}
}
