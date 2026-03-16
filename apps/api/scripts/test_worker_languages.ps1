param(
    [string]$ApiUrl = "http://localhost:3000",
    [ValidateSet("python", "node", "cpp", "java", "all")]
    [string]$Language = "all",
    [int]$TimeoutSeconds = 120,
    [int]$PollIntervalSeconds = 2,
    [switch]$Help
)

if ($Help) {
    Write-Host "Usage: .\test_worker_languages.ps1 [-ApiUrl URL] [-Language python|node|cpp|java|all] [-TimeoutSeconds 120] [-PollIntervalSeconds 2]"
    exit 0
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

$languages = @()
if ($Language -eq "all") {
    $languages = @("python", "node", "cpp", "java")
}
else {
    $languages = @($Language)
}

Write-Host "[info] Languages: $($languages -join ', ')"

try {
    $challengeBody = @{
        title = "Worker Multi-Language Smoke"
        description = "Given two integers, print their sum."
        difficulty = "easy"
        timeLimit = 2000
        memoryLimit = 256
        tags = @("worker", "e2e", "smoke")
        publicTestCases = @(
            @{
                name = "sample-1"
                input = "2 3`n"
                output = "5`n"
            }
        )
        hiddenTestCases = @(
            @{
                name = "hidden-1"
                input = "10 32`n"
                output = "42`n"
            },
            @{
                name = "hidden-2"
                input = "-5 2`n"
                output = "-3`n"
            }
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
    $testUser = "worker_test_${suffix}_${rand}"
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

    $results = @{}
    $failCount = 0

    foreach ($lang in $languages) {
        Write-Host "[info] submitting language: $lang"

        $submissionBody = @{
            challengeId = $challengeId
            code = (Get-CodeByLanguage -Lang $lang)
            language = $lang
        }

        try {
            $submission = Invoke-ApiPost -Url "$ApiUrl/submissions" -Body $submissionBody -Token $token
        }
        catch {
            Write-Error "Submission failed for $lang. $_"
            $results[$lang] = "submission_error"
            $failCount++
            continue
        }

        $submissionId = $submission.id
        if (-not $submissionId) {
            Write-Error "Submission response missing id for $lang"
            $results[$lang] = "missing_submission_id"
            $failCount++
            continue
        }

        Write-Host "[ok] queued submission: $submissionId"

        $startTime = Get-Date
        $finalStatus = "queued"

        while ($true) {
            try {
                $submissionState = Invoke-ApiGet -Url "$ApiUrl/submissions/$submissionId"
                $status = "$($submissionState.status)"
                if (-not $status) {
                    $status = "unknown"
                }
            }
            catch {
                Write-Error "Could not fetch submission $submissionId for $lang. $_"
                $status = "poll_error"
                $finalStatus = $status
                break
            }

            $elapsed = [int]((Get-Date) - $startTime).TotalSeconds
            Write-Host "[info] $lang submission=$submissionId status=$status elapsed=${elapsed}s"

            if (Is-TerminalStatus -Status $status) {
                $finalStatus = $status
                break
            }

            if ($elapsed -ge $TimeoutSeconds) {
                $finalStatus = "timeout"
                break
            }

            Start-Sleep -Seconds $PollIntervalSeconds
        }

        $results[$lang] = $finalStatus
        if ($finalStatus -ne "accepted") {
            $failCount++
        }
    }

    Write-Host ""
    Write-Host "================ Worker Language Test Summary ================"
    foreach ($lang in $languages) {
        Write-Host "${lang}: $($results[$lang])"
    }
    Write-Host "Challenge ID: $challengeId"
    Write-Host ""

    if ($failCount -gt 0) {
        Write-Error "[result] FAIL ($failCount language(s) not accepted). Ensure API, Redis, worker and runner images are up."
        exit 2
    }

    Write-Host "[result] PASS (all languages accepted)"
    exit 0
}
catch {
    Write-Error "Fatal error: $_"
    exit 1
}
