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
	"testing"
	"time"
)

var nowFunc = time.Now

func TestParseCronExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      bool
	}{
		{"*/15 0 1,15 * 1-5 /usr/bin/find", "*/15 0 1,15 * 1-5", false},
		{"@hourly", "0 * * * *", false},
		{"@reboot", "@reboot", false},
		{"invalid-expression", "", true},
	}

	for _, test := range tests {
		result, err := ParseCronExpression(test.input)

		if test.err && err == nil {
			t.Errorf("Expected error for input: %s, but got none", test.input)
		}

		if !test.err && err != nil {
			t.Errorf("Unexpected error for input: %s - %v", test.input, err)
		}

		// Check only if there's no error and the result is not empty
		if !test.err && result != "" && result != test.expected {
			t.Errorf("Expected: %s, Got: %s", test.expected, result)
		}
	}
}

func TestParseCronJob(t *testing.T) {
	tests := []struct {
		input    string
		expected CronJob
		err      bool
	}{
		{"*/15 0 1,15 * 1-5 /usr/bin/find", CronJob{"*/15", "0", "1,15", "*", "1-5", "/usr/bin/find"}, false},
		{"@reboot", CronJob{}, true}, // Updated to expect an error
		{"invalid-line", CronJob{}, true},
	}

	for _, test := range tests {
		result, err := ParseCronJob(test.input)

		if test.err && err == nil {
			t.Errorf("Expected error for input: %s, but got none", test.input)
		}

		if !test.err && err != nil {
			t.Errorf("Unexpected error for input: %s - %v", test.input, err)
		}

		if result != test.expected {
			t.Errorf("Expected: %+v, Got: %+v", test.expected, result)
		}
	}
}

func TestCronToDuration(t *testing.T) {
	now := time.Now()
	nowFunc = func() time.Time { return now }
	defer func() { nowFunc = time.Now }()

	cronJob := CronJob{"0", "12", "1", "1", "*", "/usr/bin/command1"}

	duration, err := CronToDuration(cronJob)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the calculated duration is a reasonable positive value
	if duration <= 0 {
		t.Errorf("Invalid duration: %v", duration)
	}
}
