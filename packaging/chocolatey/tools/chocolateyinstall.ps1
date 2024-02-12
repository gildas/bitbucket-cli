$ErrorActionPreference = 'Stop' # stop on all errors
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  fileType      = 'exe'
  file64        = "$toolsDir\bitbucket-cli-0.10.1-windows-amd64.7z"
  softwareName  = 'bitbucket-cli*'
  checksum64    = '4d0cf672a51800cb4409550fdaa5ce8e32440d9bb74f5c729c8866cbb4ad32d4'
  checksumType64= 'sha256'
}

Get-ChocolateyUnzip @packageArgs
Remove-Item -Path $packageArgs.file64 -Force

Write-Output "To load tab completion in your current PowerShell session, please run:"
Write-Output "  bb completion powershell | Out-String | Invoke-Expression"
Write-Output " "
Write-Output "To load completions for every new session, add the output of the above command to your powershell profile."
