import 'package:flutter_test/flutter_test.dart';
import 'package:irohasu_jaspr/components/toast.dart';

void main() {
  test('ToastController emits show and dismiss events', () async {
    final events = <ToastEvent>[];
    final subscription = ToastController.instance.events.listen(events.add);

    final id = ToastController.instance.success('Đã lưu');
    ToastController.instance.dismiss(id);
    await Future<void>.delayed(Duration.zero);

    expect(events, hasLength(2));
    expect(events.first, isA<ToastShown>());
    expect((events.first as ToastShown).toast.message, 'Đã lưu');
    expect(events.last, isA<ToastDismissed>());
    await subscription.cancel();
  });
}
