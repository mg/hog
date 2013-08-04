#!/bin/sh
# usage: markdown FILE.md | ../md.awk | pbcopy
awk '
/<p>CODE<\/p>/ { 
	print "<div class=\"code\"><pre>CODE</pre></div>" 
	next
}
/<p>.*<\/p>/ {
	printf("<span class=\"Apple-style-span\" style=\"font-size: medium;\">%s<br/><br/></span>", substr($0,4,length($0)-7))
	next
}
/^\s+$/ { 
	print "<br><br>" 
	next
}
/.+/ {
	print $0
	next
}
'