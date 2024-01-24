/*
* xim scheduler
* Copyright (C)  Alex Gaetano Padula
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package cronparser

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// CronJob represents a cron job.
type CronJob struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Command    string
}

// ParseCronExpression parses a cron expression and returns the standard cron representation.
func ParseCronExpression(expression string) (string, error) {
	components := strings.Fields(expression)

	if len(components) == 1 {
		// Check for non-standard expressions
		switch components[0] {
		case "@hourly":
			return "0 * * * *", nil
		case "@daily", "@midnight":
			return "0 0 * * *", nil
		case "@weekly":
			return "0 0 * * 0", nil
		case "@monthly":
			return "0 0 1 * *", nil
		case "@yearly", "@annually":
			return "0 0 1 1 *", nil
		case "@reboot":
			return "@reboot", nil
		}
	} else if len(components) != 6 {
		return "", fmt.Errorf("invalid crontab expression: %s, expected 6 fields", expression)
	}

	for _, component := range components {
		if component == "@reboot" {
			return "", fmt.Errorf("invalid crontab expression: @reboot should be used alone")
		}
	}

	// Ensure that there are at least 5 components before joining
	if len(components) < 5 {
		return "", fmt.Errorf("invalid crontab expression: %s", expression)
	}

	return strings.Join(components[:5], " "), nil
}

// ParseCronJob parses a cron line and returns a CronJob struct.
func ParseCronJob(cronLine string) (CronJob, error) {
	fields := strings.Fields(cronLine)

	if len(fields) < 6 {
		return CronJob{}, fmt.Errorf("invalid crontab line: %s", cronLine)
	}

	return CronJob{
		Minute:     fields[0],
		Hour:       fields[1],
		DayOfMonth: fields[2],
		Month:      fields[3],
		DayOfWeek:  fields[4],
		Command:    strings.Join(fields[5:], " "),
	}, nil
}

// CronToDuration calculates the duration until the next scheduled time for a CronJob.
func CronToDuration(cron CronJob) (time.Duration, error) {
	// Convert cron fields to time layout
	//layout := fmt.Sprintf("%s %s %s %s *", cron.Minute, cron.Hour, cron.DayOfMonth, cron.Month)

	// Parse the layout to get the next scheduled time
	nextScheduledTime, err := GetNextScheduledTime(cron)
	if err != nil {
		return 0, err
	}

	// Calculate the duration until the next scheduled time
	durationUntilNext := nextScheduledTime.Sub(time.Now())

	return durationUntilNext, nil
}

// GetNextScheduledTime finds the next valid time based on the cron fields.
func GetNextScheduledTime(cron CronJob) (time.Time, error) {
	now := time.Now()

	// Parse the cron expression fields
	minuteValues, err := parseCronField(cron.Minute, 0, 59)
	if err != nil {
		return time.Time{}, err
	}

	hourValues, err := parseCronField(cron.Hour, 0, 23)
	if err != nil {
		return time.Time{}, err
	}

	dayOfMonthValues, err := parseCronField(cron.DayOfMonth, 1, 31)
	if err != nil {
		return time.Time{}, err
	}

	monthValues, err := parseCronField(cron.Month, 1, 12)
	if err != nil {
		return time.Time{}, err
	}

	// Find the next valid time based on the cron fields
	for {
		now = now.Add(time.Minute)

		if contains(minuteValues, now.Minute()) &&
			contains(hourValues, now.Hour()) &&
			contains(dayOfMonthValues, now.Day()) &&
			contains(monthValues, int(now.Month())) {
			break
		}
	}

	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location()), nil
}

// Contains checks if a value is present in a slice.
func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicate values from a slice.
func removeDuplicates(arr []int) []int {
	var result []int
	seen := make(map[int]bool)

	for _, num := range arr {
		if !seen[num] {
			result = append(result, num)
			seen[num] = true
		}
	}

	return result
}

// GetLastDayOfMonth returns the last day of the current month.
func GetLastDayOfMonth() int {
	now := time.Now()
	lastDayOfMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location())
	return lastDayOfMonth.Day()
}

// ParseCronField parses a cron field and returns the corresponding values.
func parseCronField(field string, minValue, maxValue int) ([]int, error) {
	if field == "*" {
		// Wildcard represents all values
		return allValuesInRange(minValue, maxValue), nil
	}

	if field == "L" {
		// "L" character represents the last day of the month
		return []int{GetLastDayOfMonth()}, nil
	}

	if strings.HasPrefix(field, "*/") {
		step, err := strconv.Atoi(strings.TrimPrefix(field, "*/"))
		if err != nil {
			return nil, err
		}
		return stepValuesInRange(minValue, maxValue, step), nil
	}

	if strings.Contains(field, ",") {
		// Handle comma-separated values, e.g., "1,15"
		parts := strings.Split(field, ",")
		var result []int

		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}

			if num < minValue || num > maxValue {
				return nil, fmt.Errorf("value out of range in cron field: %s", field)
			}

			result = append(result, num)
		}

		// Sort and remove duplicates
		sort.Ints(result)
		result = removeDuplicates(result)

		return result, nil
	}

	// Handle single numeric value
	num, err := strconv.Atoi(field)
	if err != nil {
		return nil, err
	}

	if num < minValue || num > maxValue {
		return nil, fmt.Errorf("value out of range in cron field: %s", field)
	}

	return []int{num}, nil
}

// AllValuesInRange returns all values in the specified range.
func allValuesInRange(min, max int) []int {
	var result []int
	for i := min; i <= max; i++ {
		result = append(result, i)
	}
	return result
}

// StepValuesInRange returns values in the specified range with a given step.
func stepValuesInRange(min, max, step int) []int {
	var result []int
	for i := min; i <= max; i += step {
		result = append(result, i)
	}
	return result
}

// FindNextValue finds the next valid value based on the current value.
func findNextValue(values []int, current int) int {
	for _, value := range values {
		if value >= current {
			return value
		}
	}
	// If no value is found, return the first value in the list
	return values[0]
}

// PrintCronJobInfo prints information about a CronJob.
func PrintCronJobInfo(crontab []string) {
	for i, entry := range crontab {
		parsed, err := ParseCronJob(entry)
		if err != nil {
			log.Fatalf("Error parsing cron job #%d: %v", i+1, err)
		}

		cronExpression, err := ParseCronExpression(entry)
		if err != nil {
			log.Fatalf("Error parsing cron expression #%d: %v", i+1, err)
		}

		duration, err := CronToDuration(parsed)
		if err != nil {
			log.Fatalf("Error calculating duration for cron job #%d: %v", i+1, err)
		}

		fmt.Printf("Cron Job #%d: %s\n", i+1, cronExpression)
		fmt.Printf("Parsed Result: %+v\n", parsed)
		fmt.Printf("Next Run Duration: %v\n", duration)
		fmt.Println()
	}
}
