## xim

xim is a multi-platform scheduler that supports standard cron expressions with six fields: Minute, Hour, Day of Month, Month, Day of Week, and Command.

xim handles non-standard expressions such as "@hourly," "@daily," "@weekly," "@monthly," "@yearly," "@annually," and "@reboot."

### ✨ Features
⭐ Provides detailed error handling for invalid cron expressions, ensuring that the parser returns meaningful error messages for easier debugging.

⭐ Calculates the duration until the next scheduled execution time for a given cron job.

⭐ Dynamically determines the last day of the month when encountering the "L" character in the DayOfMonth field.

⭐ Supports step values in cron fields, such as "*/5" in the minute field, indicating every 5 minutes.

⭐ Handles comma-separated values in cron fields, allowing multiple specific values to be specified (e.g., "1,15" in the day of month field).

* Supports the use of "*" as a wildcard character, representing all values in a given field.

⭐ Removes duplicate values when encountering comma-separated values in cron fields.

⭐ Calculates the next scheduled time efficiently, considering the specific cron expression fields.

⭐ Allows for extension or modification based on specific use cases or additional features.

### ximtab
is the ximtable that contains all your jobs.
```
* * * * * echo hello world
*/15 0 1,12 * 1-5 echo hello world again
```