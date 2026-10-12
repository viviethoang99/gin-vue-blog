const _apiBaseUrl = String.fromEnvironment('API_BASE_URL');
const _publicApiBaseUrl = String.fromEnvironment('PUBLIC_API_BASE_URL');

/// Internal URL used by the Jaspr server to call Gin.
final String apiBaseUrl = _normalizeApiBaseUrl(_apiBaseUrl);

/// Browser-visible URL used by search and article assets.
///
/// In production this is normally `/`, because Nginx exposes Gin under the
/// same origin. It may be an absolute URL during local development.
final String publicApiBaseUrl = _normalizePublicBaseUrl(
  _publicApiBaseUrl.isEmpty ? _apiBaseUrl : _publicApiBaseUrl,
);

String _normalizeApiBaseUrl(String value) {
  final normalized = value.trim();
  final uri = Uri.tryParse(normalized);
  if (normalized.isEmpty || uri == null || !uri.hasAuthority || (uri.scheme != 'http' && uri.scheme != 'https')) {
    throw StateError(
      'API_BASE_URL must be a valid HTTP(S) URL. '
      'Load it with --dart-define-from-file=.env.',
    );
  }
  return normalized.endsWith('/') ? normalized : '$normalized/';
}

String _normalizePublicBaseUrl(String value) {
  final normalized = value.trim();
  if (normalized.isEmpty) {
    throw StateError(
      'PUBLIC_API_BASE_URL must be an HTTP(S) URL or an absolute path.',
    );
  }
  final uri = Uri.tryParse(normalized);
  if (uri == null ||
      (uri.hasScheme && uri.scheme != 'http' && uri.scheme != 'https') ||
      (!uri.hasScheme && !normalized.startsWith('/'))) {
    throw StateError(
      'PUBLIC_API_BASE_URL must be an HTTP(S) URL or an absolute path.',
    );
  }
  return normalized.endsWith('/') ? normalized : '$normalized/';
}
