Import-Module powershell-yaml
$pythonScriptPath = "C:\Users\rogla\bachelor\go-incentive-simulation\main.go"
$outputFilePath = "C:\Users\rogla\bachelor\go-incentive-simulation\output.txt"

# Set the number of times to run the script
$numRuns = 50
$endGr = 26
$startGr = 6

#rm C:\Users\rogla\bachelor\go-incentive-simulation\output.txt

# Loop through the number of runs and run the script each time
for ($i = $startGr; $i -le $endGr; $i++) {
    $i >> $outputFilePath
    $data = Get-Content -Path "C:\Users\rogla\bachelor\go-incentive-simulation\config.yaml" -Raw | ConvertFrom-Yaml
    $data.Confoptions.NumGoroutines = $i
    $data | ConvertTo-Yaml | Set-Content -Path "C:\Users\rogla\bachelor\go-incentive-simulation\config.yaml"
    for ($j = 1; $j -le $numRuns; $j++) {
        # Set the output file path with the current iteration number
        # Run the Python script and redirect its output to the output file
        go run $pythonScriptPath >> $outputFilePath
    }
}
