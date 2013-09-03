#!/bin/sh
# usage: markdown FILE.md | ../md.awk | pbcopy
awk '
NR < 3 {
	next
}
/<p>CODE<\/p>/ { 
	print "<div class=\"code\"><pre>CODE</pre></div>\n" 
	next
}
/<\/code><\/pre>/ { 
	print "</pre></div><br/>\n" 
	next
}
/<p>.*<\/p>/ {
	printf("<span class=\"Apple-style-span\" style=\"font-size: medium;\">%s<br/><br/></span>\n\n", substr($0,4,length($0)-7))
	next
}
/^\s+$/ { 
	print "<br><br>\n\n" 
	next
}
/.+/ {
	print $0
	next
}
'