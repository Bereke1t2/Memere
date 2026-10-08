import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Mirkuz (ምርኩዝ) — Minimal Exam Preparation & Learning Platform",
  description:
    "Master the Ethiopian Grade 12 National University Entrance Examination with Mirkuz. Comprehensive video courses, downloadable offline notes, timed national mock exams, and seamless Telebirr payments.",
  keywords: [
    "Ethiopian Grade 12 Exam",
    "University Entrance Exam Ethiopia",
    "Mirkuz",
    "ምርኩዝ",
    "Matric Prep Ethiopia",
    "Memere",
    "Grade 12 Entrance Preparation",
    "Ethiopian Matric Exam",
    "Telebirr Education",
    "Natural Science Grade 12",
    "Social Science Grade 12",
  ],
  authors: [{ name: "Mirkuz Education Platform" }],
  openGraph: {
    title: "Mirkuz (ምርኩዝ) — Minimal Exam Preparation Platform",
    description:
      "Prepare for your Grade 12 University Entrance Exam with HD video lessons, offline PDF notes, and timed mock exams. Download the app on Google Play.",
    url: "https://mirkuz.app",
    siteName: "Mirkuz Exam Prep",
    locale: "en_US",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Mirkuz — Grade 12 University Entrance Prep",
    description: "Ethiopia's premier Grade 12 exam preparation app. 100% offline ready.",
  },
  icons: {
    icon: [
      { url: "/favicon.ico" },
      { url: "/favicon-32x32.png", sizes: "32x32", type: "image/png" },
      { url: "/favicon-16x16.png", sizes: "16x16", type: "image/png" },
      { url: "/icon.png", type: "image/png" },
    ],
    apple: [
      { url: "/apple-touch-icon.png", sizes: "180x180", type: "image/png" },
    ],
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-white text-[#64748B] antialiased selection:bg-[#00B894]/20 selection:text-[#1E2B58]">
        {children}
      </body>
    </html>
  );
}
