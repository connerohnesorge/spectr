# Tasks

## Edge Case Testing

- [x] 1.1 Normal task with "quotes" and 'single quotes'
- [x] 1.2 Task with backslashes: C:\Users\test\path and \\network\share
- [x] 1.3 Task with newlines: This is line 1
And this should also be part of the description (this won't parse as single line)
- [x] 1.4 Task with special JSON chars: { "key": "value" } and [array]
- [x] 1.5 Task with unicode: 🚀 Emoji test with 中文 Chinese chars and Ñoño
- [x] 1.6 Task with tabs: indented with tabs here
- [x] 1.7 Task with escape sequences: \n \t \r \b \f \" \\ \/ \u0041
- [x] 1.8 Very long description that exceeds normal length limits to test if the
  system can handle descriptions that go on and on and on with lots of text
  that might cause buffer issues or memory problems or other edge cases that
  only appear with extremely long strings that contain multiple sentences and
  potentially hundreds of characters that need to be properly escaped and
  validated during the JSON marshaling process which should handle this
  gracefully without any errors or truncation or corruption of the data
  regardless of how verbose the task description becomes over time
- [x] 1.9 Mixed special chars: C:\path\to\file.txt with "quotes" and {json}
  and \u003chtml\u003e and 50% discount
- [x] 1.10 Backslash at end: path\to\directory\
- [x] 1.11 Control chars: ASCII control characters like  (backspace) and  (form feed)
- [x] 1.12 JSON edge case: Task with "},{ which might break naive parsers
- [x] 1.13 Unicode edge: Combining diacritics é café naïve Zürich
- [x] 1.14 Math symbols: ∑ ∫ √ ± × ÷ ≠ ≈ ≤ ≥
- [x] 1.15 Currency and symbols: $100 €50 £30 ¥1000 © ® ™ § ¶
