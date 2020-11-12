package vm

func BuildP2PKHScript(pubHash []byte) []byte {
	var script []byte
	script = append(script, byte(OpDup))
	script = append(script, byte(OpHash256))
	script = append(script, pubHash...)
	script = append(script, byte(OpEqualVerify))
	script = append(script, byte(OpCheckSig))
	return script
}
