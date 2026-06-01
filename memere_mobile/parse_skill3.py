import re
import os

with open('docs/skills/phase1/SKILL_3.md', 'r') as f:
    content = f.read()

# Pattern to extract filename and code content
pattern = r'### FILE \w+ — `([^`]+)`\n\n```dart\n(.*?)```'
matches = re.findall(pattern, content, re.DOTALL)

for filename, code in matches:
    # ensure directory exists
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    with open(filename, 'w') as f:
        f.write(code)
    print(f'Wrote {filename}')
