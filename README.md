# Basic LSP

## Introduction

This is a basic LSP integration just for fun.

## Testing

This should also work on VS \*\*\*\*
This should also work on VS Code

autocompleted NeoVim (BTW)

diagnostic test for VS Code
VS Code

## Add LSP to NeoVim (BTW)

Add to `init.lua`

```
local client = vim.lsp.start_client {
  name = "basiclsp",
  cmd = {<path>/basic-lsp/main"},
  on_attach = require("lazyvim.util").lsp.on_attach,
 }
if not client then
  vim.notify "Hey, you did not setup the basiclsp client well"
  return
end
vim.api.nvim_create_autocmd("FileType", {
  pattern="markdown",
  callback = function()
    vim.lsp.buf_attach_client(0, client)
  end,
})

```

# Credits

- Credits to tjdevries
