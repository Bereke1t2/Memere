lines = open('pubspec.yaml').readlines()
with open('pubspec.yaml', 'w') as f:
    skip = False
    for line in lines:
        if line.strip() == 'fonts:':
            skip = True
        elif skip and not line.startswith(' ') and line.strip() != '':
            skip = False
        
        if not skip:
            f.write(line)
