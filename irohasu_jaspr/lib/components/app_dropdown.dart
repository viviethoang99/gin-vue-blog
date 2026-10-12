import 'package:jaspr/dom.dart';
import 'package:jaspr/dom.dart' as dom;
import 'package:jaspr/jaspr.dart';

class DropdownItem {
  const DropdownItem({required this.value, required this.label, this.disabled = false});

  final String value;
  final String label;
  final bool disabled;
}

/// Accessible native dropdown. Place it inside an `@client` component when
/// using [onChanged]; static forms can omit the callback and submit [name].
class AppDropdown extends StatelessComponent {
  const AppDropdown({
    required this.id,
    required this.items,
    this.label,
    this.name,
    this.value,
    this.placeholder,
    this.helpText,
    this.errorText,
    this.disabled = false,
    this.required = false,
    this.onChanged,
    super.key,
  });

  final String id;
  final List<DropdownItem> items;
  final String? label;
  final String? name;
  final String? value;
  final String? placeholder;
  final String? helpText;
  final String? errorText;
  final bool disabled;
  final bool required;
  final ValueChanged<String?>? onChanged;

  @override
  Component build(BuildContext context) {
    final descriptionId = errorText != null ? '$id-error' : (helpText != null ? '$id-help' : null);

    return div(classes: 'app-dropdown${errorText != null ? ' app-dropdown--error' : ''}', [
      if (label != null)
        dom.label(htmlFor: id, classes: 'app-dropdown__label', [
          .text(label!),
          if (required) span(classes: 'app-dropdown__required', [.text(' *')]),
        ]),
      div(classes: 'app-dropdown__control', [
        select(
          id: id,
          name: name,
          value: value,
          disabled: disabled,
          required: required,
          classes: 'app-dropdown__select',
          attributes: {
            'aria-describedby': ?descriptionId,
            if (errorText != null) 'aria-invalid': 'true',
          },
          onChange: onChanged == null ? null : (values) => onChanged!(values.firstOrNull),
          [
            if (placeholder != null)
              option(value: '', selected: value == null, disabled: required, [.text(placeholder!)]),
            for (final item in items)
              option(value: item.value, selected: item.value == value, disabled: item.disabled, [.text(item.label)]),
          ],
        ),
        span(classes: 'app-dropdown__chevron', attributes: const {'aria-hidden': 'true'}, [.text('⌄')]),
      ]),
      if (errorText != null)
        span(id: '$id-error', classes: 'app-dropdown__message app-dropdown__message--error', [.text(errorText!)])
      else if (helpText != null)
        span(id: '$id-help', classes: 'app-dropdown__message', [.text(helpText!)]),
    ]);
  }

  @css
  static List<StyleRule> get styles => [
    css('.app-dropdown').styles(display: .flex, flexDirection: .column, gap: .all(6.px)),
    css('.app-dropdown__label').styles(
      color: const Color.variable('--color-text-primary'),
      fontSize: 0.875.rem,
      fontWeight: .w700,
    ),
    css('.app-dropdown__required').styles(color: const Color.variable('--color-danger')),
    css('.app-dropdown__control').styles(position: .relative()),
    css('.app-dropdown__select', [
      css('&').styles(
        width: 100.percent,
        padding: .symmetric(vertical: 11.px, horizontal: 14.px),
        border: .all(style: .solid, color: const Color.variable('--color-border'), width: 1.px),
        radius: .all(.circular(10.px)),
        outline: const Outline(style: .none),
        appearance: .none,
        cursor: .pointer,
        color: const Color.variable('--color-text-primary'),
        fontSize: 0.95.rem,
        backgroundColor: const Color.variable('--color-surface'),
        raw: const {'padding-right': '42px'},
      ),
      css('&:focus-visible').styles(
        border: .all(style: .solid, color: const Color.variable('--color-focus'), width: 1.px),
        raw: const {'box-shadow': '0 0 0 3px rgba(22, 137, 216, 0.2)'},
      ),
      css('&:disabled').styles(opacity: 0.55, cursor: .notAllowed),
    ]),
    css('.app-dropdown__chevron').styles(
      position: .absolute(top: 50.percent, right: 14.px),
      pointerEvents: .none,
      transform: .translate(y: (-50).percent),
      color: const Color.variable('--color-text-secondary'),
    ),
    css('.app-dropdown__message').styles(color: const Color.variable('--color-text-secondary'), fontSize: 0.75.rem),
    css('.app-dropdown__message--error').styles(color: const Color.variable('--color-danger')),
    css('.app-dropdown--error .app-dropdown__select').styles(
      border: .all(style: .solid, color: const Color.variable('--color-danger'), width: 1.px),
    ),
  ];
}
