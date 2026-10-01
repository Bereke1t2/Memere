import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Mirkuz (ምርኩዝ) — Ethiopian Grade 12 University Entrance Exam Prep",
  description:
    "Master the Ethiopian Grade 12 National University Entrance Examination with Mirkuz. Comprehensive video courses, downloadable offline notes, timed national mock exams, and seamless Telebirr payments.",
  keywords: [
    "Ethiopian Grade 12 Exam",
    "University Entrance Exam Ethiopia",
    "Mirkuz",
    "Memere",
    "Grade 12 Entrance Preparation",
    "Ethiopian Matric Exam",
    "Telebirr Education",
    "Ethiopian High School Online Learning",
    "Natural Science Grade 12",
    "Social Science Grade 12",
  ],
  authors: [{ name: "Mirkuz Education Platform" }],
  openGraph: {
    title: "Mirkuz (ምርኩዝ) — Grade 12 National Exam Prep",
    description:
      "Prepare for your Grade 12 University Entrance Exam with HD video lessons, offline PDF notes, and timed mock exams. Download the app on Google Play.",
    url: "https://mirkuz.et",
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
    icon: "/favicon.ico",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className="min-h-screen bg-[#050505] text-[#f4f4f5] antialiased selection:bg-[#10b981]/30 selection:text-white">
        {children}
      </body>
    </html>
  );
}
