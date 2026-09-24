import Flutter
import UIKit

@main
@objc class AppDelegate: FlutterAppDelegate, FlutterImplicitEngineDelegate {
  override func application(
    _ application: UIApplication,
    didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
  ) -> Bool {
    // 1. Listen for screen recording / mirroring state changes
    NotificationCenter.default.addObserver(
      self,
      selector: #selector(screenCaptureChanged),
      name: UIScreen.capturedDidChangeNotification,
      nil
    )

    // 2. Add secure subview layer to protect against screenshots
    if let window = self.window {
      makeWindowSecure(window)
    }

    return super.application(application, didFinishLaunchingWithOptions: launchOptions)
  }

  func didInitializeImplicitFlutterEngine(_ engineBridge: FlutterImplicitEngineBridge) {
    GeneratedPluginRegistrant.register(with: engineBridge.pluginRegistry)
  }

  // 3. Obscure app snapshot in iOS App Switcher
  override func applicationWillResignActive(_ application: UIApplication) {
    super.applicationWillResignActive(application)
    guard let window = self.window else { return }
    if window.viewWithTag(9999) != nil { return }
    let blur = UIBlurEffect(style: .dark)
    let blurView = UIVisualEffectView(effect: blur)
    blurView.frame = window.bounds
    blurView.tag = 9999
    window.addSubview(blurView)
  }

  override func applicationDidBecomeActive(_ application: UIApplication) {
    super.applicationDidBecomeActive(application)
    self.window?.viewWithTag(9999)?.removeFromSuperview()
  }

  @objc private func screenCaptureChanged() {
    if UIScreen.main.isCaptured {
      showScreenRecordingShield()
    } else {
      removeScreenRecordingShield()
    }
  }

  private func showScreenRecordingShield() {
    guard let window = self.window else { return }
    if window.viewWithTag(8888) != nil { return }
    let shieldView = UIView(frame: window.bounds)
    shieldView.backgroundColor = UIColor(red: 0.04, green: 0.07, blue: 0.09, alpha: 1.0)
    shieldView.tag = 8888

    let label = UILabel()
    label.text = "Screen Recording Protected"
    label.textColor = .white
    label.font = UIFont.systemFont(ofSize: 18, weight: .bold)
    label.textAlignment = .center
    label.translatesAutoresizingMaskIntoConstraints = false
    shieldView.addSubview(label)

    NSLayoutConstraint.activate([
      label.centerXAnchor.constraint(equalTo: shieldView.centerXAnchor),
      label.centerYAnchor.constraint(equalTo: shieldView.centerYAnchor)
    ])

    window.addSubview(shieldView)
  }

  private func removeScreenRecordingShield() {
    self.window?.viewWithTag(8888)?.removeFromSuperview()
  }

  private func makeWindowSecure(_ window: UIWindow) {
    let field = UITextField()
    field.isSecureTextEntry = true
    window.addSubview(field)
    field.centerYAnchor.constraint(equalTo: window.centerYAnchor).isActive = true
    field.centerXAnchor.constraint(equalTo: window.centerXAnchor).isActive = true
    window.layer.superlayer?.addSublayer(field.layer)
    field.layer.sublayers?.first?.addSublayer(window.layer)
  }
}

