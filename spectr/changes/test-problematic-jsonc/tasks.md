# Tasks

## Edge Case Testing

- [ ] 1.1 Normal task with "quotes" and 'single quotes'
- [ ] 1.2 Task with backslashes: C:\Users\test\path and \\network\share
- [ ] 1.3 Task with newlines: This is line 1
And this should also be part of the description (this won't parse as single line)
- [ ] 1.4 Task with special JSON chars: { "key": "value" } and [array]
- [ ] 1.5 Task with unicode: 🚀 Emoji test with 中文 Chinese chars and Ñoño
- [ ] 1.6 Task with tabs:	indented	with	tabs	here
- [ ] 1.7 Task with escape sequences: \n \t \r \b \f \" \\ \/ \u0041
- [ ] 1.8 Very long description that exceeds normal length limits to test if the system can handle descriptions that go on and on and on with lots of text that might cause buffer issues or memory problems or other edge cases that only appear with extremely long strings that contain multiple sentences and potentially hundreds of characters that need to be properly escaped and validated during the JSON marshaling process which should handle this gracefully without any errors or truncation or corruption of the data regardless of how verbose the task description becomes over time
- [ ] 1.9 Mixed special chars: C:\path\to\file.txt with "quotes" and {json} and <html> and 50% discount
- [ ] 1.10 Backslash at end: path\to\directory\
- [ ] 1.11 Control chars: ASCII control characters like  (backspace) and  (form feed)
- [ ] 1.12 JSON edge case: Task with "},{ which might break naive parsers
- [ ] 1.13 Unicode edge: Combining diacritics é café naïve Zürich
- [ ] 1.14 Math symbols: ∑ ∫ √ ± × ÷ ≠ ≈ ≤ ≥
- [ ] 1.15 Currency and symbols: $100 €50 £30 ¥1000 © ® ™ § ¶
