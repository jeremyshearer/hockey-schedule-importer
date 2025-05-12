# Hockey Schedule Importer

Hockey Schedule Importer is a utility script that converts a csv of games exported from Gamesheets to the import format required by benchapp

## Running the program

1. Download a csv of your schedule from gamesheets (in this example named `gamesheet_schedule.csv`)
2. run the importer:

```go
go run main.go -in gamesheet_schedule.csv -out benchapp_schedule.csv
```

3. Upload benchapp_schedule.csv to Bench App
