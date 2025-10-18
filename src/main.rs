mod frontend;
mod utils;

use crate::frontend::speak_sphere::SpeakSphere;
use crate::utils::display_size::get_unit_size;
use crate::utils::font::setup_fonts;

use eframe::egui;

fn main() -> Result<(), eframe::Error> {
    let (unit_width, unit_height) = get_unit_size();
    let (w, h) = (unit_width * 50.0, unit_height * 50.0);
    let options = eframe::NativeOptions {
        viewport: egui::ViewportBuilder::default()
            .with_inner_size([w, h])
            .with_title("Speak Sphere"),
        ..Default::default()
    };
    eframe::run_native(
        "Speak Sphere",
        options,
        Box::new(|cc| {
            setup_fonts(&cc.egui_ctx);
            Ok(Box::new(SpeakSphere::new()))
        }),
    )
}
