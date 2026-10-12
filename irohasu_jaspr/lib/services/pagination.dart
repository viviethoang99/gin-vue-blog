import '../models/paginated_result.dart';

PaginatedResult<T> paginateItems<T>(List<T> items, {required int page, required int pageSize}) {
  if (pageSize <= 0) {
    throw ArgumentError.value(pageSize, 'pageSize', 'Must be greater than zero');
  }

  final totalPages = items.isEmpty ? 1 : (items.length / pageSize).ceil();
  final currentPage = page.clamp(1, totalPages);
  final start = (currentPage - 1) * pageSize;
  final end = (start + pageSize).clamp(0, items.length);

  return PaginatedResult<T>(
    items: items.sublist(start, end),
    currentPage: currentPage,
    pageSize: pageSize,
    totalItems: items.length,
    totalPages: totalPages,
  );
}
