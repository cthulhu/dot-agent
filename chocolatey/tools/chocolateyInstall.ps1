$ErrorActionPreference = 'Stop';

$packageName = 'dot-agent'
$toolsDir    = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$version     = $env:ChocolateyPackageVersion

if ([string]::IsNullOrEmpty($version)) {
  $version = '0.7.0'
}

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url           = "https://github.com/cthulhu/dot-agent/releases/download/v$version/dot-agent_${version}_windows_amd64.zip"
  softwareName  = 'dot-agent'
  checksum      = 'e42231cfdaf7d6192302077b43e0aff463f83af10c3f2fced8a1ca4116b577e0'
  checksumType  = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs

# Rename the extracted dot-agent_0.7.0_windows_amd64.exe to dot-agent.exe so that Chocolatey automatically shims it.
$extractedExe = Get-ChildItem -Path $toolsDir -Filter "dot-agent_*.exe" | Select-Object -First 1
if ($extractedExe) {
  Rename-Item -Path $extractedExe.FullName -NewName "dot-agent.exe"
}
