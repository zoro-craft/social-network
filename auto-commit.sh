#!/bin/bash

# كيقرا ݣاع الملفات لي فيهم تغييرات (modified, deleted, untracked)
git status --porcelain | while read -r line; do
  # كيعزل المسار ديال الملف والحالة ديالو
  file="${line:3}"
  status="${line:0:2}"
  
  # كياخد غير السمية ديال الملف و الدوسي ديالو باش يصاوب ميساج زوين
  base=$(basename "$file")
  dir=$(basename $(dirname "$file"))
  
  # كيحيد الإمتدادات (.js, .css, .module) وكيرد كلشي حروف صغيرة
  clean_base=$(echo "$base" | sed -E 's/\.module//g; s/\.(js|jsx|css)//g' | tr '[:upper:]' '[:lower:]')
  clean_dir=$(echo "$dir" | tr '[:upper:]' '[:lower:]')

  # كيصاوب الميساج على حساب شنو وقع فـ الملف (ماكسيموم 4 كلمات)
  if [[ "$status" == *"D"* ]]; then
    msg="remove $clean_base"
  elif [[ "$status" == "??" ]]; then
    msg="add $clean_base styles"
  else
    if [[ "$base" == "page.js" ]]; then
      msg="update $clean_dir page"
    elif [[ "$base" == "globals.css" ]]; then
      msg="tweak global css"
    else
      msg="fix $clean_base"
    fi
  fi

  # كيدير add و commit
  git add "$file"
  git commit -m "$msg"
  
  # كيطبع الميساج باش تعرف فين وصل، وكيتسنى دقيقة
  echo "✅ committed: $msg | waiting 1 min..."
  sleep 60
done

echo "🎉 all files committed successfully!"