# Tasks

## Extreme Edge Cases

- [x] 1.1 Nested quotes: "He said \"Hello\" and she replied \"Hi there\""
- [x] 1.2 Multiple backslashes: C:\\\\\\\\path\\\\\\\\to\\\\\\\\file
- [x] 1.3 JSON injection attempt: "},{"id":"injected","status":"hacked
- [x] 1.4 Unicode edge: null bytes (if any), zero-width chars: ​‌‍, RTL marks: ‏‎
- [x] 1.5 HTML injection: <script>alert('xss')</script> and <!-- comment -->
- [x] 1.6 Path traversal: ../../../etc/passwd and ..\..\..\..\windows\system32
- [x] 1.7 Quote bombardment: """""""""""""""""""""""""""""""""""""""
- [x] 1.8 Backslash bombardment: \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
- [x] 1.9 Mixed bombardment: \"\"\"\"\\ \"\"\"\"\\ \"\"\"\"\\ \"\"\"\"\\
- [x] 1.10 JSONC comment injection: // comment */ or /* comment */ in description
- [x] 1.11 Literal newline char test: Line1\nLine2\nLine3 (with literal \n not actual newline)
- [x] 1.12 Empty-ish: "" and '' and ``
- [x] 1.13 Format string: %s %d %x %n ${var} #{var} {{var}}
- [x] 1.14 Regex chars: .* .+ .? ^$ [a-z] (foo|bar) \d+ \w+ \s+
- [x] 1.15 All printable ASCII: !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_`abcdefghijklmnopqrstuvwxyz{|}~
