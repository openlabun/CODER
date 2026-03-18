param(
    [string]$ApiUrl = "http://localhost:3000",
    [ValidateSet("python", "node", "cpp", "java")]
    [string]$Language = "python",
    [int]$SubmissionCount = 5,
    [int]$TimeoutSeconds = 180,
    [int]$PollIntervalSeconds = 2,
    [switch]$Help
)

if ($Help) {
    Write-Host "Usage: .\test_parallel_queue_5.ps1 [-ApiUrl URL] [-Language python|node|cpp|java] [-SubmissionCount 5] [-TimeoutSeconds 180] [-PollIntervalSeconds 2]"
    exit 0
}

if ($SubmissionCount -lt 1) {
    Write-Error "SubmissionCount must be >= 1"
    exit 1
}

$ErrorActionPreference = "Stop"

function Invoke-ApiPost {
    param(
        [string]$Url,
        [object]$Body,
        [string]$Token
    )

    $headers = @{}
    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }

    return Invoke-RestMethod -Method Post -Uri $Url -Headers $headers -ContentType "application/json" -Body ($Body | ConvertTo-Json -Depth 10)
}

function Invoke-ApiGet {
    param(
        [string]$Url,
        [string]$Token
    )

    $headers = @{}
    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }

    return Invoke-RestMethod -Method Get -Uri $Url -Headers $headers
}

function Get-CodeByLanguage {
    param([string]$Lang)

    switch ($Lang) {
        "python" {
            return @"
import sys
nums = list(map(int, sys.stdin.read().split()))
print(sum(nums[:2]))
"@
        }
        "node" {
            return @"
const fs = require('fs');
const nums = fs.readFileSync(0, 'utf8').trim().split(/\s+/).map(Number);
console.log((nums[0] || 0) + (nums[1] || 0));
"@
        }
        "cpp" {
            return @"
#include <bits/stdc++.h>
using namespace std;
int main() {
    long long a, b;
    if (!(cin >> a >> b)) return 0;
    cout << (a + b) << "\n";
    return 0;
}
"@
        }
        "java" {
            return @"
import java.util.*;
public class Main {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        if (!sc.hasNextLong()) return;
        long a = sc.nextLong();
        long b = sc.hasNextLong() ? sc.nextLong() : 0;
        System.out.println(a + b);
    }
}
"@
        }
        default {
            throw "Unsupported language: $Lang"
        }
    }
}

function Is-TerminalStatus {
    param([string]$Status)
    return @("accepted", "wrong_answer", "error") -contains $Status
}

Write-Host "[info] API URL: $ApiUrl"
Write-Host "[info] Language: $Language"
Write-Host "[info] Submission burst size: $SubmissionCount"

try {
    $challengeBody = @{
        title = "Parallel Queue Burst Test"
        description = "Given two integers, print their sum."
        difficulty = "easy"
        timeLimit = 2000
        memoryLimit = 256
        tags = @("queue", "parallel", "burst")
        publicTestCases = @(
            @{ name = "sample-1"; input = "2 3`n"; output = "5`n" }
        )
        hiddenTestCases = @(
            @{ name = "hidden-1"; input = "10 32`n"; output = "42`n" },
            @{ name = "hidden-2"; input = "100 200`n"; output = "300`n" }
        )
    }

    $challenge = Invoke-ApiPost -Url "$ApiUrl/challenges" -Body $challengeBody
    $challengeId = $challenge.id
    if (-not $challengeId) {
        throw "Challenge creation did not return an id"
    }
    Write-Host "[ok] challenge created: $challengeId"

    $suffix = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
    $rand = Get-Random -Minimum 1000 -Maximum 9999
    $testUser = "queue_test_${suffix}_${rand}"
    $testPass = "pass_${suffix}_${rand}"

    $registerBody = @{
        username = $testUser
        password = $testPass
        role = "student"
    }

    $registerResp = Invoke-ApiPost -Url "$ApiUrl/auth/register" -Body $registerBody
    $token = $registerResp.accessToken
    if (-not $token) {
        throw "Register response did not return accessToken"
    }
    Write-Host "[ok] student registered: $testUser"

    # Enqueue burst submissions as fast as possible.
    $sourceCode = Get-CodeByLanguage -Lang $Language
    $submissionRecords = @()
    $enqueueStart = Get-Date

    for ($i = 1; $i -le $SubmissionCount; $i++) {
        $body = @{
            challengeId = $challengeId
            code = $sourceCode
            language = $Language
        }

        $resp = Invoke-ApiPost -Url "$ApiUrl/submissions" -Body $body -Token $token
        if (-not $resp.id) {
            throw "Submission #$i did not return id"
        }

        $submissionRecords += [PSCustomObject]@{
            Index = $i
            Id = $resp.id
            Status = "queued"
            LastStatus = "queued"
            StartTime = Get-Date
            EndTime = $null
            ElapsedSeconds = 0
        }

        Write-Host "[ok] queued #$i -> $($resp.id)"
    }

    $enqueueElapsed = [math]::Round(((Get-Date) - $enqueueStart).TotalMilliseconds, 0)
    Write-Host "[info] enqueue burst completed in ${enqueueElapsed}ms"

    # Poll all submissions until all reach terminal state or timeout.
    $globalStart = Get-Date
    while ($true) {
        $allDone = $true
        $globalElapsed = [int]((Get-Date) - $globalStart).TotalSeconds

        foreach ($record in $submissionRecords) {
            if (Is-TerminalStatus -Status $record.Status) {
                continue
            }

            if ($record.Status -eq "timeout" -or $record.Status -eq "poll_error") {
                continue
            }

            $allDone = $false
            try {
                $state = Invoke-ApiGet -Url "$ApiUrl/submissions/$($record.Id)"
                $current = "$($state.status)"
                if (-not $current) { $current = "unknown" }
                $record.Status = $current
                $record.LastStatus = $current
            }
            catch {
                $record.Status = "poll_error"
                $record.LastStatus = "poll_error"
            }

            $elapsed = [int]((Get-Date) - $record.StartTime).TotalSeconds
            $record.ElapsedSeconds = $elapsed

            if (Is-TerminalStatus -Status $record.Status -or $record.Status -eq "poll_error") {
                $record.EndTime = Get-Date
            }
            elseif ($elapsed -ge $TimeoutSeconds -or $globalElapsed -ge $TimeoutSeconds) {
                $record.Status = "timeout"
                $record.EndTime = Get-Date
            }
        }

        $statusLine = ($submissionRecords | ForEach-Object { "#$($_.Index)=$($_.Status)" }) -join " | "
        Write-Host "[info] poll t=${globalElapsed}s :: $statusLine"

        $unfinished = $submissionRecords | Where-Object {
            -not (Is-TerminalStatus -Status $_.Status) -and $_.Status -ne "timeout" -and $_.Status -ne "poll_error"
        }

        if ($allDone -or $unfinished.Count -eq 0) {
            break
        }

        if ($globalElapsed -ge $TimeoutSeconds) {
            break
        }

        Start-Sleep -Seconds $PollIntervalSeconds
    }

    $accepted = ($submissionRecords | Where-Object { $_.Status -eq "accepted" }).Count
    $wrongAnswer = ($submissionRecords | Where-Object { $_.Status -eq "wrong_answer" }).Count
    $errorCount = ($submissionRecords | Where-Object { $_.Status -eq "error" }).Count
    $timeoutCount = ($submissionRecords | Where-Object { $_.Status -eq "timeout" }).Count
    $pollErrCount = ($submissionRecords | Where-Object { $_.Status -eq "poll_error" }).Count

    Write-Host ""
    Write-Host "================ Parallel Queue Summary ================"
    foreach ($record in $submissionRecords) {
        Write-Host "#${($record.Index)} id=$($record.Id) status=$($record.Status) elapsed=${($record.ElapsedSeconds)}s"
    }
    Write-Host ""
    Write-Host "accepted=$accepted wrong_answer=$wrongAnswer error=$errorCount timeout=$timeoutCount poll_error=$pollErrCount"
    Write-Host "challengeId=$challengeId"

    if ($accepted -ne $SubmissionCount) {
        Write-Error "[result] FAIL (expected $SubmissionCount accepted, got $accepted)."
        exit 2
    }

    Write-Host "[result] PASS (all $SubmissionCount submissions accepted)"
    exit 0
}
catch {
    Write-Error "Fatal error: $_"
    exit 1
}
