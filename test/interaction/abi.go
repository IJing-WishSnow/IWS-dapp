package interaction

import _ "embed"

// ==================== 编译时嵌入 ABI 文件 ====================

//go:embed contracts/token/ERC20.abi
var erc20ABIJSON string
