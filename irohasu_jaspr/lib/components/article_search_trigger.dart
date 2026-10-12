import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';

import 'article_search.dart';

@client
class ArticleSearchTrigger extends StatelessComponent {
  const ArticleSearchTrigger({this.classes, super.key});

  final String? classes;

  @override
  Component build(BuildContext context) {
    return button(
      type: .button,
      classes: ['article-search-trigger', ?classes].join(' '),
      attributes: const {'aria-label': 'Tìm kiếm bài viết'},
      onClick: ArticleSearchController.instance.open,
      [
        span(
          classes: 'material-symbols-outlined',
          attributes: const {'aria-hidden': 'true'},
          [.text('search')],
        ),
      ],
    );
  }
}
