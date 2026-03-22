package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jeremyshearer/hockey-schedule-importer/converter"
	"github.com/spf13/cobra"
)

var addr string

var httpdCmd = &cobra.Command{
	Use:   "httpd",
	Short: "Start an HTTP server that converts posted CSV",
	RunE: func(cmd *cobra.Command, args []string) error {
		http.HandleFunc("/", handleForm)
		http.HandleFunc("/convert", handleConvert)
		fmt.Println("Listening on", addr)
		return http.ListenAndServe(addr, nil)
	},
}

const formHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Hockey Schedule Converter</title>
<script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="min-h-screen bg-gray-100 flex items-center justify-center">
<div class="bg-white shadow-md rounded-lg p-8 max-w-lg w-full">
  <h1 class="text-2xl font-bold text-gray-800 mb-6">Hockey Schedule Converter</h1>

  <div class="mb-6 bg-gray-50 border border-gray-200 rounded-lg p-4">
    <h2 class="text-sm font-semibold text-gray-700 mb-3">Step 1: Export from GameSheet</h2>
    <ol class="text-sm text-gray-600 list-decimal list-inside space-y-1 mb-4">
      <li>Go to your team page on GameSheet (e.g. <a href="https://gamesheetstats.com/seasons/14401/teams/493736/schedule" class="text-blue-600 hover:underline">HatTrick Swayzes</a>)</li>
      <li>Click the <span class="font-medium">Schedule</span> tab</li>
      <li>Click the <span class="inline-flex items-center font-semibold text-xs bg-gray-800 text-white px-2 py-0.5 rounded">EXPORT</span> button in the bottom-right corner</li>
    </ol>
    <h2 class="text-sm font-semibold text-gray-700 mb-3">Step 2: Convert</h2>
    <p class="text-sm text-gray-600 mb-4">Upload the exported CSV below and click Convert to download the BenchApp-compatible file.</p>
    <h2 class="text-sm font-semibold text-gray-700 mb-3">Step 3: Import into BenchApp</h2>
    <ol class="text-sm text-gray-600 list-decimal list-inside space-y-1">
      <li>In BenchApp, click <span class="font-medium">Add</span></li>
      <li>Choose <span class="font-medium">Import Schedule</span> from the dropdown</li>
      <li>Upload the converted CSV file</li>
    </ol>
  </div>

  <form method="POST" action="/convert" enctype="multipart/form-data" class="space-y-4">
    <label class="block">
      <span class="text-sm font-medium text-gray-700">Upload CSV</span>
      <input type="file" name="file" accept=".csv" required
        class="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100">
    </label>
    <button type="submit"
      class="w-full bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 font-medium">
      Convert
    </button>
  </form>
</div>
</body>
</html>`

func handleForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, formHTML)
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

	var csvBytes []byte
	var err error

	file, _, fileErr := r.FormFile("file")
	if fileErr == nil {
		defer file.Close()
		csvBytes, err = converter.Convert(file, os.Stderr)
	} else {
		csvBytes, err = converter.Convert(r.Body, os.Stderr)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="schedule.csv"`)
	w.Write(csvBytes)
}

func init() {
	httpdCmd.Flags().StringVar(&addr, "addr", ":8080", "Listen address")
	Root.AddCommand(httpdCmd)
}
