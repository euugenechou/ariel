package misc

import "bytes"

func Flounder(msg string) string {
	var out bytes.Buffer
	out.WriteString("\"" + msg + "\"\n")
	out.WriteString("    \\\n")
	out.WriteString("      \\\n")
	out.WriteString("        , __\n")
	out.WriteString("        \\`\\\"._     _,\n")
	out.WriteString("        / _  |||;._/ )\n")
	out.WriteString("      _/@ @  ///  | (\n")
	out.WriteString("     ( (`__,     ,`\\|\n")
	out.WriteString("      '.\\_/ |\\_.'\n")
	out.WriteString("        `\"\"```")
	return out.String()
}
