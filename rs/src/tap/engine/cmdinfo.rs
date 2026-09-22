use crate::engine::structs::{CmdInfo};

impl<'a, H> CmdInfo<'a, H> {
    pub fn new(
        handler: H,
        required: Vec<&'a str>,
        optional: Vec<&'a str>,
        unlimited: bool,
        docs: &'static str,
    ) -> Self {
        Self {
            docs: docs,
            handler,
            required,
            optional,
            unlimited,
        }
    }
}
