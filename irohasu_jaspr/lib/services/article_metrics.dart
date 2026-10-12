String plainArticleText(String value) {
  return value
      .replaceAll(RegExp(r'<[^>]*>'), ' ')
      .replaceAll(RegExp(r'!\[[^\]]*\]\([^)]*\)'), ' ')
      .replaceAll(RegExp(r'\[([^\]]+)\]\([^)]*\)'), r'$1')
      .replaceAll(RegExp(r'[`#>*_~|\-]'), ' ')
      .replaceAll(RegExp(r'''[.,!?;:…“”"'‘’{}]'''), ' ')
      .replaceAll(RegExp(r'\s+'), ' ')
      .trim();
}

int countArticleWords(String value) {
  final plainText = plainArticleText(value);
  if (plainText.isEmpty) return 0;
  return plainText.split(' ').length;
}

String formatVietnameseInteger(int value) {
  final digits = value.toString();
  final buffer = StringBuffer();
  for (var index = 0; index < digits.length; index++) {
    if (index > 0 && (digits.length - index) % 3 == 0) buffer.write('.');
    buffer.write(digits[index]);
  }
  return buffer.toString();
}
