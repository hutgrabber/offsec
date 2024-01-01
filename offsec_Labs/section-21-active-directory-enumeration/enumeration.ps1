$DomainObj = [System.DirectoryServices.ActiveDirectory.Domain]::GetCurrentDomain();
# Primary Domain Controller ($PDC):
$PDC = $DomainObj.PdcRoleOwner.Name;
$PDC;
