$ErrorActionPreference = 'Stop';

$packageName = 'dot-agent'
$toolsDir    = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$version     = $env:ChocolateyPackageVersion

if ([string]::IsNullOrEmpty($version)) {
  $version = '0.8.0'
}

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url           = "https://github.com/cthulhu/dot-agent/releases/download/v$version/dot-agent_${version}_windows_amd64.zip"
  softwareName  = 'dot-agent'
  checksum      = 'cf2c39aec155e322cffbcef204ed17f174aabb0a4e01a209f31e01008a6e294c'
  checksumType  = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs

# Rename the extracted dot-agent_0.8.0_windows_amd64.exe to dot-agent.exe so that Chocolatey automatically shims it.
$extractedExe = Get-ChildItem -Path $toolsDir -Filter "dot-agent_*.exe" | Select-Object -First 1
if ($extractedExe) {
  Rename-Item -Path $extractedExe.FullName -NewName "dot-agent.exe"
}
