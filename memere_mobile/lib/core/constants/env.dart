import 'package:envied/envied.dart';

part 'env.g.dart';

@Envied(path: '.env')
abstract class Env {
  @EnviedField(varName: 'BASE_URL')
  static const String baseUrl = _Env.baseUrl;

  @EnviedField(varName: 'CHAPA_PUBLIC_KEY', obfuscate: true)
  static final String chapaPublicKey = _Env.chapaPublicKey;

  @EnviedField(varName: 'FIREBASE_PROJECT_ID')
  static const String firebaseProjectId = _Env.firebaseProjectId;
}
