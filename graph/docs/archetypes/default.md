---
title: "{{ replace .TranslationBaseName "-" " " | title }}"
date: {{ .Date }}
anchor: "{{ replace .TranslationBaseName "-" " " | title | urlize }}"
weight:
---
