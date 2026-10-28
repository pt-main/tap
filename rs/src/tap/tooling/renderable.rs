use std::{collections::HashMap, error::Error, fmt, sync::Arc};

use string_template::Template;

use crate::tooling::colors::set::colorize;

type RenderScopeT<'a> = HashMap<&'a str, &'a str>;
type RendererT<'a> = Arc<dyn Fn(&RenderScopeT<'a>) -> Result<String, Box<dyn Error>> + 'a>;

#[derive(Debug)]
pub enum RenderError {
    RendererNotFound(String),
    RendererError(Box<dyn Error>),
    InvalidValue {
        key: String,
        kind: &'static str,
        value: String,
    },
}

impl fmt::Display for RenderError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::RendererNotFound(name) => write!(f, "renderer not found: {name}"),
            Self::RendererError(err) => write!(f, "renderer error: {err}"),
            Self::InvalidValue { key, kind, value } => {
                write!(f, "{key:?}:{kind} has invalid value: {value:?}")
            }
        }
    }
}

impl Error for RenderError {}

#[derive(Default, Clone)]
pub struct Renderable<'a> {
    pub renderers: HashMap<String, RendererT<'a>>,
    pub scope: RenderScopeT<'a>,
}

impl<'a> fmt::Debug for Renderable<'a> {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.debug_struct("Renderable")
            .field("renderers", &self.renderers.keys().collect::<Vec<_>>())
            .field("scope", &self.scope)
            .finish()
    }
}

#[derive(Clone, Copy)]
enum Kind {
    Str,
    Int,
    Bool,
    Raw,
}

fn split_kind(key: &str) -> (&str, Kind) {
    if let Some(k) = key.strip_suffix(":int") {
        (k, Kind::Int)
    } else if let Some(k) = key.strip_suffix(":bool") {
        (k, Kind::Bool)
    } else if let Some(k) = key.strip_suffix(":raw") {
        (k, Kind::Raw)
    } else {
        (key, Kind::Str)
    }
}

impl<'a> Renderable<'a> {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn render_json(&self) -> Result<String, RenderError> {
        let mut items: Vec<(&str, &str, Kind)> = self
            .scope
            .iter()
            .map(|(k, v)| {
                let (name, kind) = split_kind(k);
                (name, *v, kind)
            })
            .collect();
        items.sort_by_key(|(name, _, _)| *name);

        let mut out = String::from("{");
        for (i, (name, value, kind)) in items.into_iter().enumerate() {
            if i > 0 {
                out.push(',');
            }
            push_json_string(&mut out, name);
            out.push(':');
            match kind {
                Kind::Str => push_json_string(&mut out, value),
                Kind::Int => {
                    if value.parse::<f64>().is_err() {
                        return Err(RenderError::InvalidValue {
                            key: name.to_string(),
                            kind: "int",
                            value: value.to_string(),
                        });
                    }
                    out.push_str(value);
                }
                Kind::Bool => {
                    if value != "true" && value != "false" {
                        return Err(RenderError::InvalidValue {
                            key: name.to_string(),
                            kind: "bool",
                            value: value.to_string(),
                        });
                    }
                    out.push_str(value);
                }
                Kind::Raw => out.push_str(value),
            }
        }
        out.push('}');
        Ok(out)
    }

    pub fn render_with(&self, name: &str) -> Result<String, RenderError> {
        let renderer = self
            .renderers
            .get(name)
            .ok_or_else(|| RenderError::RendererNotFound(name.to_string()))?;

        let mut scope: RenderScopeT<'a> = HashMap::with_capacity(self.scope.len());
        for (k, v) in &self.scope {
            let (key, _) = split_kind(k);
            scope.insert(key, v);
        }

        renderer(&scope).map_err(RenderError::RendererError)
    }
}

pub fn render_json_array(rows: &[Renderable<'_>]) -> Result<String, RenderError> {
    let parts: Vec<String> = rows
        .iter()
        .map(Renderable::render_json)
        .collect::<Result<_, _>>()?;
    Ok(format!("[{}]", parts.join(",")))
}

pub fn add_template_render<'a>(
    templ: &'a str,
    name: impl Into<String>,
    renderable: &mut Renderable<'a>,
) {
    let template = Template::new(templ);
    let func: RendererT<'a> = Arc::new(
        move |scope: &RenderScopeT<'a>| -> Result<String, Box<dyn Error>> {
            Ok(template.render(scope))
        },
    );
    renderable.renderers.insert(name.into(), func);
}

pub fn add_template_render_colored<'a>(
    templ: &'a str,
    name: impl Into<String>,
    renderable: &mut Renderable<'a>,
) {
    let template = Template::new(templ);
    let func: RendererT<'a> = Arc::new(
        move |scope: &RenderScopeT<'a>| -> Result<String, Box<dyn Error>> {
            Ok(colorize(&template.render(scope)))
        },
    );
    renderable.renderers.insert(name.into(), func);
}

fn push_json_string(out: &mut String, s: &str) {
    use std::fmt::Write;
    out.push('"');
    for c in s.chars() {
        match c {
            '"' => out.push_str("\\\""),
            '\\' => out.push_str("\\\\"),
            '\n' => out.push_str("\\n"),
            '\r' => out.push_str("\\r"),
            '\t' => out.push_str("\\t"),
            '\u{08}' => out.push_str("\\b"),
            '\u{0C}' => out.push_str("\\f"),
            c if (c as u32) < 0x20 => {
                let _ = write!(out, "\\u{:04x}", c as u32);
            }
            c => out.push(c),
        }
    }
    out.push('"');
}
