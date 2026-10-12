import 'package:flutter_test/flutter_test.dart';
import 'package:irohasu_jaspr/services/pagination.dart';

void main() {
  group('paginateItems', () {
    final items = List<int>.generate(12, (index) => index + 1);

    test('returns the requested page', () {
      final result = paginateItems(items, page: 2, pageSize: 5);

      expect(result.items, [6, 7, 8, 9, 10]);
      expect(result.currentPage, 2);
      expect(result.totalPages, 3);
      expect(result.hasPreviousPage, isTrue);
      expect(result.hasNextPage, isTrue);
    });

    test('returns the remaining items on the last page', () {
      final result = paginateItems(items, page: 3, pageSize: 5);

      expect(result.items, [11, 12]);
      expect(result.hasNextPage, isFalse);
    });

    test('clamps page numbers outside the available range', () {
      expect(paginateItems(items, page: 0, pageSize: 5).currentPage, 1);
      expect(paginateItems(items, page: 99, pageSize: 5).currentPage, 3);
    });

    test('rejects a non-positive page size', () {
      expect(() => paginateItems(items, page: 1, pageSize: 0), throwsArgumentError);
    });
  });
}
