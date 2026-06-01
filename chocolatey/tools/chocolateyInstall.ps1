$ErrorActionPreference = 'Stop';

$packageName = 'dot-agent'
$toolsDir    = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$version     = $env:ChocolateyPackageVersion

if ([string]::IsNullOrEmpty($version)) {
  $version = '0.8.3'
}

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url           = "https://github.com/cthulhu/dot-agent/releases/download/v$version/dot-agent_${version}_windows_amd64.zip"
  softwareName  = 'dot-agent'
  checksum      = '' # TODO: Update after release
  checksumType  = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs

# Rename the extracted dot-agent_0.8.3_windows_amd64.exe to dot-agent.exe so that Chocolatey automatically shims it.
$extractedExe = Get-ChildItem -Path $toolsDir -Filter "dot-agent_*.exe" | Select-Object -First 1
if ($extractedExe) {
  Rename-Item -Path $extractedExe.FullName -NewName "dot-agent.exe"
}
