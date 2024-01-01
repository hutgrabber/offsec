$PDC = [System.DirectoryServices.ActiveDirectory.Domain]::GetCurrentDomain().PdcRoleOwner.Name

$DN = ([adsi]'').distinguishedName 

$LDAP = "LDAP://$PDC/$DN"

$direntry = New-Object System.DirectoryServices.DirectoryEntry($LDAP);

$searcher = New-Object System.DirectoryServices.DirectorySearcher($direntry);

$searcher.FindAll()
