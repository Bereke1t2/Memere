import React from "react";
import Link from "next/link";
import { GraduationCap, ArrowLeft, BookOpen, FileCheck, AlertCircle, ShieldAlert } from "lucide-react";

export const metadata = {
  title: "Terms of Service — Mirkuz (ምርኩዝ)",
  description: "Terms and conditions of use for the Mirkuz Grade 12 University Entrance Learning Platform.",
};

export default function TermsPage() {
  return (
    <div className="min-h-screen bg-[#0a0a0c] text-[#e2e8f0]">
      {/* Top Header */}
      <header className="border-b border-[#1e293b] bg-[#0f172a]/70 backdrop-blur sticky top-0 z-50">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 py-4 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-[#10b981] to-[#059669] flex items-center justify-center text-white font-bold">
              <GraduationCap className="w-5 h-5" />
            </div>
            <div>
              <span className="text-lg font-bold tracking-tight text-white">Mirkuz</span>
              <span className="ml-2 px-1.5 py-0.5 text-[10px] font-semibold bg-[#10b981]/15 text-[#10b981] rounded border border-[#10b981]/30">
                ምርኩዝ
              </span>
            </div>
          </Link>
          <Link
            href="/"
            className="inline-flex items-center gap-2 text-xs font-semibold text-[#94a3b8] hover:text-white transition-colors bg-[#1e293b] px-3.5 py-2 rounded-lg border border-[#334155]"
          >
            <ArrowLeft className="w-3.5 h-3.5" />
            Back to Home
          </Link>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-4xl mx-auto px-4 sm:px-6 py-12 space-y-10">
        {/* Title banner */}
        <div className="border-b border-[#1e293b] pb-8 space-y-3">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-[#38bdf8]/10 border border-[#38bdf8]/25 text-[#38bdf8] text-xs font-semibold">
            <BookOpen className="w-4 h-4" />
            Platform Service Agreement
          </div>
          <h1 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
            Terms of Service
          </h1>
          <p className="text-sm text-[#94a3b8]">
            Effective Date: October 2026 • Mirkuz Education Platform
          </p>
        </div>

        {/* Section 1 */}
        <section className="space-y-3">
          <h2 className="text-lg font-bold text-white flex items-center gap-2">
            <span className="text-[#10b981]">1.</span> Acceptance of Terms
          </h2>
          <p className="text-xs sm:text-sm text-[#cbd5e1] leading-relaxed">
            By creating an account or accessing Mirkuz (including mobile applications, web portals, course materials, video lectures, and mock exams), you agree to be bound by these Terms of Service. If you do not agree with these terms, you should not use the platform.
          </p>
        </section>

        {/* Section 2 */}
        <section className="space-y-3">
          <h2 className="text-lg font-bold text-white flex items-center gap-2">
            <span className="text-[#10b981]">2.</span> Educational Purpose & License
          </h2>
          <p className="text-xs sm:text-sm text-[#cbd5e1] leading-relaxed">
            All curriculum videos, question banks, study guides, and mock exams are protected intellectual property created for individual academic preparation for Ethiopian Grade 12 students. You are granted a limited, personal, non-exclusive, non-transferable license to access the content.
          </p>
          <div className="p-4 rounded-xl bg-[#0f172a] border border-[#1e293b] space-y-2">
            <div className="flex items-center gap-2 text-xs font-bold text-[#f59e0b]">
              <AlertCircle className="w-4 h-4" />
              Prohibited Content Usage:
            </div>
            <ul className="list-disc list-inside text-xs text-[#94a3b8] space-y-1 pl-2">
              <li>Recording, ripping, re-uploading, or publicly broadcasting proprietary video lessons.</li>
              <li>Scraping exam question banks or distributing answer keys commercially.</li>
              <li>Sharing or renting account credentials across multiple individuals.</li>
            </ul>
          </div>
        </section>

        {/* Section 3 */}
        <section className="space-y-3">
          <h2 className="text-lg font-bold text-white flex items-center gap-2">
            <span className="text-[#10b981]">3.</span> Payments & Access Subscriptions
          </h2>
          <p className="text-xs sm:text-sm text-[#cbd5e1] leading-relaxed">
            Premium course access, full national mock entrance exams, and chapter downloads can be unlocked via supported payment gateways (e.g. Telebirr, Chapa). Access remains valid for the duration of the active subscription or individual course enrollment tier.
          </p>
        </section>

        {/* Section 4 */}
        <section className="space-y-3">
          <h2 className="text-lg font-bold text-white flex items-center gap-2">
            <span className="text-[#10b981]">4.</span> Account Termination & Data Deletion
          </h2>
          <p className="text-xs sm:text-sm text-[#cbd5e1] leading-relaxed">
            You may terminate your account at any time via the mobile app profile settings or our{" "}
            <Link href="/delete-account" className="text-[#38bdf8] underline">
              Account Deletion Page
            </Link>. We reserve the right to suspend or terminate accounts that engage in harassment, unauthorized content redistribution, or fraudulent activity.
          </p>
        </section>

        {/* Section 5 */}
        <section className="space-y-3 border-t border-[#1e293b] pt-8">
          <h2 className="text-lg font-bold text-white flex items-center gap-2">
            <FileCheck className="w-5 h-5 text-[#10b981]" />
            5. Governing Law
          </h2>
          <p className="text-xs text-[#94a3b8] leading-relaxed">
            These terms are governed by the laws and regulations of the Federal Democratic Republic of Ethiopia.
          </p>
          <p className="text-xs text-[#94a3b8]">
            Questions? Contact <a href="mailto:support@mirkuz.et" className="text-[#38bdf8] underline">support@mirkuz.et</a>.
          </p>
        </section>
      </main>

      {/* Footer */}
      <footer className="border-t border-[#1e293b] py-8 text-center text-xs text-[#64748b]">
        <p>© {new Date().getFullYear()} Mirkuz (ምርኩዝ) Education Platform. All rights reserved.</p>
      </footer>
    </div>
  );
}
