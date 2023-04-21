Import-Module powershell-yaml
$goProgram = "C:\Users\rogla\bachelor\go-incentive-simulation\main.go"
$outputFilePath = "C:\Users\rogla\bachelor\go-incentive-simulation\output.txt"

$numRuns = 50
$endGr = 26
$startGr = 6

rm C:\Users\rogla\bachelor\go-incentive-simulation\output.txt

for ($i = $endGr; $i -le $startGr; $i++) {
    $i >> $outputFilePath
    $data = Get-Content -Path "C:\Users\rogla\bachelor\go-incentive-simulation\config.yaml" -Raw | ConvertFrom-Yaml
    $data.Confoptions.NumGoroutines = $i
    $data | ConvertTo-Yaml | Set-Content -Path "C:\Users\rogla\bachelor\go-incentive-simulation\config.yaml"
    for ($j = 1; $j -le $numRuns; $j++) {
        go run $goProgram >> $outputFilePath
    }
}
