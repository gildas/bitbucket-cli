$ErrorActionPreference = 'Stop' # stop on all errors
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$isArm64    = ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64' -or $env:PROCESSOR_ARCHITEW6432 -eq 'ARM64')

if ($isArm64) {
  $file     = "$toolDir\bitbucket-cli-{{VERSION}}-windows-arm64.7z"
  $checksum = '{{CHECKSUM_ARM64}}'
} else {
  $file     = "$toolDir\bitbucket-cli-{{VERSION}}-windows-amd64.7z"
  $checksum = '{{CHECKSUM_AMD64}}'
}

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  fileType      = 'exe'
  file64        = $file
  softwareName  = 'bitbucket-cli*'
  checksum64    = $checksum
  checksumType64= 'sha256'
}

Get-ChocolateyUnzip @packageArgs
Remove-Item -Path $packageArgs.file64 -Force

Write-Output "To load tab completion in your current PowerShell session, please run:"
Write-Output "  bb completion powershell | Out-String | Invoke-Expression"
Write-Output " "
Write-Output "To load completions for every new session, add the output of the above command to your powershell profile."
