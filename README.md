### What is Cron scheduler ?

A cron scheduler is a time-based job scheduler. It allows users to schedule scripts or commands to run automatically at specified times or intervals. It runs in the background and executes scheduled tasks according to the configuration in a file called `crontab`.

Here's a quick overview of how it works:

1. **Crontab File**: This file contains a list of jobs and their schedules. Each line in the crontab file represents a job and follows a specific syntax to define when the job should run.
    
2. **Scheduling Syntax**: The schedule for each job is defined using a combination of fields:
    
    * Minute (0-59)
        
    * Hour (0-23)
        
    * Day of the month (1-31)
        
    * Month (1-12)
        
    * Day of the week (0-7, where both 0 and 7 represent Sunday)
        
    
    For example, a line in a crontab file might look like this:
    
    ```bash
    30 8 * * 1-5 /path/to/script.sh
    ```
    
    This would run `/path/to/`[`script.sh`](http://script.sh) at 8:30 AM every Monday through Friday.
    
3. **Cron Jobs**: These are the actual commands or scripts specified in the crontab file. Each job is executed according to its scheduled time.
    

Cron is widely used for tasks such as automating backups, sending emails, and running periodic maintenance scripts.

### What is Orbis ?

Orbis is a robust and lightweight cron scheduler, crafted in the Go programming language. It features an API server that serves as an interface to the scheduler. Orbis leverages Postgres for storing cron jobs and employs a resizable goroutine pool to achieve efficient concurrency.

To interpret scheduling expressions, Orbis utilizes the `cronexpr` library (available at [github.com/gorhill/cronexpr](http://github.com/gorhill/cronexpr)[). This library enables O](https://github.com/gorhill/cronexpr)rbis to accurately determine when commands should be activated based on the provided schedule.

### Features of Orbis

* **Concurrency**: Leveraging Go’s goroutines, Orbis efficiently handles multiple jobs simultaneously.
    
* **Flexibility**: Supports standard cron exp[ressions for flexible sched](https://github.com/gorhill/cronexpr)uling.
    
* **Simplicity**: Minimal setup required with clear separation of concerns between scheduling and execution.
    

### Key libraries used

* Cronexpr - [https://github.com/gorhill/cronexpr](https://github.com/gorhill/cronexpr)
    
* Gin - [https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
    

### Getting Started with Orbis

Orbis is built on top of Go’s powerful concurrency model and integrates seamlessly with popular libraries such as Gin for web routing and Cronexpr for cron expression parsing.

#### Step 1

Clone the repository -&gt; [https://github.com/heyyakash/orbis](https://github.com/heyyakash/orbis)

#### Step 2

go to the root folder and create a new file called .env

```bash
touch .env
```

#### Step 3

Populate the .env file with following information (using the same is fine for testing)

```bash
POSTGRES_PASSWORD = mysecret
POSTGRES_USER = postgres
POSTGRES_DB = postgres
POSTGRES_PORT = 5432
POSTGRES_HOST = pg-image
POOL = 5
```

* `POSTGRES_PASSWORD` : refers to the password for postgres database
    
* `POSTGRES_USER` : refers to the user of the postgres database
    
* `POSTGRES_DB` : refers to the database that will store the table
    
* `POSTGRES_PORT` : refers to the port at which the postgres db will run
    
* `POSTGRES_HOST` : refers to the host address for the PostgreSQL database (usually the name of the PostgreSQL container).
    
* `POOL` : refers to the size of goroutine pool
    

#### Step 4

Run the following command to start the scheduler

```bash
docker compose up
```

#### Step 5

Check whether the scheduler is running or not by running the following (assuming the scheduler is running on the port 8080)

```bash
curl localhost:8080/ping
```

This should return the following result

```json
{"message":"pong"}
```

### Interacting with the API Server (Examples)

#### 1\. Add a job

POST `/set`

```bash
curl -X POST \
  'http://localhost:8080/set' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "command":"echo Akash Sharma",
  "schedule":"*/1 * * * *"
}'
```

Returns the id of the job created

```json
{
  "message": "Added job with id : 3"
}
```

#### 2\. Get a job

GET : `/:id`

```bash
curl  -X GET 'http://localhost:8080/13'
```

Returns the job with the given id

```json
{
  "JobId": 3,
  "Command": "curl localhost:8080/ping",
  "Schedule": "*/1 * * * *",
  "NextRun": "2024-07-22T09:25:00Z"
}
```

#### 3\. Get all jobs

GET `/all`

```bash
curl  -X GET \
  'http://localhost:8080/all'
```

returns all jobs

```json
{
  "data": [
    {
      "JobId": 1,
      "Command": "echo Akash D",
      "Schedule": "*/1 * * * *",
      "NextRun": "2024-07-22T09:26:00Z"
    },
    {
      "JobId": 2,
      "Command": "curl localhost:8080/ping",
      "Schedule": "*/1 * * * *",
      "NextRun": "2024-07-22T09:26:00Z"
    },
    {
      "JobId": 3,
      "Command": "curl localhost:8080/ping",
      "Schedule": "*/1 * * * *",
      "NextRun": "2024-07-22T09:26:00Z"
    }
  ]
}
```

#### 4\. Delete a job

DELETE `/:id`

```bash
curl -X DELETE 'http://localhost:8080/11'
```

#### 5\. Delete all jobs

DELETE `/all`

```bash
curl  -X DELETE 'http://localhost:8080/all' 
```

### Future Enhancements

* Switching to another lightweight database
    
* Add Metrics and logging
    
* Retry mechanisms
