# Load .env file and set environment variables
Get-Content "configs/.env" | ForEach-Object {
    if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
        $name = $matches[1].Trim()
        $value = $matches[2].Trim()
        [System.Environment]::SetEnvironmentVariable($name, $value)
    }
}

# Read CONNECTION_STRING from environment
$connectionString = $env:CONNECTION_STRING

if (-not $connectionString) {
    Write-Error "❌ CONNECTION_STRING not found in .env file"
    exit 1
}

# Run migration using the env var
Write-Host "Running migrations with connection: $connectionString"
migrate -path db/migrations -database $connectionString up
