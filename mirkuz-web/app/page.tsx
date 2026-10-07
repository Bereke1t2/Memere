import React from "react";
import { Navbar } from "@/components/Navbar";
import { HeroSection } from "@/components/HeroSection";
import { ClientsBar } from "@/components/ClientsBar";
import { CommunitySection } from "@/components/CommunitySection";
import { QuoteBanner } from "@/components/QuoteBanner";
import { HowWeDoIt } from "@/components/HowWeDoIt";
import { FeatureSplitTwo } from "@/components/FeatureSplitTwo";
import { AppShowcase } from "@/components/AppShowcase";
import { TestimonialSpotlight } from "@/components/TestimonialSpotlight";
import { BlogMarketingSection } from "@/components/BlogMarketingSection";
import { InteractiveExamDemo } from "@/components/InteractiveExamDemo";
import { StreamCurriculum } from "@/components/StreamCurriculum";
import { PricingSection } from "@/components/PricingSection";
import { FaqSection } from "@/components/FaqSection";
import { Footer } from "@/components/Footer";

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-white text-[#64748B] selection:bg-[#00B894]/20 selection:text-[#1E2B58] relative">
      {/* 1. Header / Navbar */}
      <Navbar />

      {/* Main Page Flow matching Dribbble Showcase Design Language */}
      <main>
        {/* 2. Hero Section: Two-tone headline, Deep Navy button, SVG puzzle/rocket artwork */}
        <HeroSection />

        {/* 3. Partner Universities & Payment Gateways */}
        <ClientsBar />

        {/* 4. We Focus On: 3-Pillar Cards (Curriculum, Strategy, Offline) */}
        <CommunitySection />

        {/* 5. Wide Glassmorphic Quote Banner */}
        <QuoteBanner />

        {/* 6. How We Do It: 4-Step Horizontal Timeline (01-04) + Navy Banner 1 */}
        <HowWeDoIt />

        {/* 7. Strategic 2x2 Feature Grid & Isometric Dashboard Art */}
        <FeatureSplitTwo />

        {/* 8. Midnight Navy Mobile App Showcase & Real App Screenshots */}
        <AppShowcase />

        {/* 9. What Students Say: Star Rating Testimonials + Navy Banner 2 */}
        <TestimonialSpotlight />

        {/* 10. National Exam Insights & Study Guides */}
        <BlogMarketingSection />

        {/* 11. Interactive Exam Simulator Micro-Experience */}
        <InteractiveExamDemo />

        {/* 12. Comprehensive Stream Curriculum (Natural & Social) */}
        <StreamCurriculum />

        {/* 13. Transparent ETB Pricing & Local Payment Gateways */}
        <PricingSection />

        {/* 14. Frequently Asked Questions */}
        <FaqSection />
      </main>

      {/* 15. Midnight Navy Footer */}
      <Footer />
    </div>
  );
}
