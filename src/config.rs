use rdev::display_size;

pub(crate) fn get_display_size() -> (u64, u64) {
    let (w, h) = display_size().expect("Failed to get display size");
    (w, h)
}
