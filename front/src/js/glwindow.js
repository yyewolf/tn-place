// ***************************************************
// ***************************************************
// GLSL code for the vertex shader
// Scales and rotates the quad
//
const viewportVertexShaderSource = `
	precision mediump float;
	attribute vec2 vert;
	uniform vec2 cam;
	uniform vec2 tex_scale;
	uniform vec2 view_scale;
	uniform float zoom;
	varying vec2 uv;
	void main() {
		uv = vert + 0.5;
		vec2 pos = ((vert * tex_scale - cam) * zoom) / view_scale;
		pos += 0.5;
		pos.y = 1.0 - pos.y;
		gl_Position = vec4(pos * 2.0 - 1.0, 0.0, 1.0);
	}
`;

// ***************************************************
// ***************************************************
// GLSL code for the fragment shader
// Paints the texture onto the quad
//
const viewportFragmentShaderSource = `
	precision mediump float;
	uniform sampler2D tex;
	uniform vec2 tex_size;
	varying vec2 uv;
	uniform vec2 mouse;
	void main() {
		vec4 pixel_color = texture2D(tex, uv);
		if (pixel_color.a < 1.0) {
			vec2 pixel_coord = fract(vec2(uv.x * tex_size.x, uv.y * tex_size.y));
      vec4 other_color = vec4(abs(0.8-pixel_color.x), abs(0.8-pixel_color.y), abs(0.8-pixel_color.z), 1.0);
			if (pixel_coord.x <= 0.08 && pixel_coord.y <= 0.3 || pixel_coord.y <= 0.08 && pixel_coord.x <= 0.3) {
        pixel_color = other_color;
      } else if (pixel_coord.x >= 0.92 && pixel_coord.y <= 0.3 || pixel_coord.y <= 0.08 && pixel_coord.x >= 0.7) {
        pixel_color = other_color;
      } else if (pixel_coord.x <= 0.08 && pixel_coord.y >= 0.7 || pixel_coord.y >= 0.92 && pixel_coord.x <= 0.3) {
        pixel_color = other_color;
      } else if (pixel_coord.x >= 0.92 && pixel_coord.y >= 0.7 || pixel_coord.y >= 0.92 && pixel_coord.x >= 0.7) {
        pixel_color = other_color;
      }
		}
		gl_FragColor = pixel_color;
	}
`;

export class GLWindow {
  #cvs;
  #gl;
  #program;
  #tex;
  #texFramebuffer;
  #texScale;
  #camPos;
  #zoom;

  #u_cam;
  #u_zoom;
  #u_tex;
  #u_view;
  #a_vert;
  #u_size;

  outline = { x: 0, y: 0, originalColor: new Uint8Array([0, 0, 0]) };

  constructor(cvs) {
    this.#cvs = cvs;
    this.#gl = cvs.getContext("webgl");
    if (this.#gl == null) {
      return;
    }

    this.#texScale = { x: 0, y: 0 };
    this.#camPos = { x: 0, y: 0 };
    this.#zoom = 1;

    const vertexShader = this.#compileShader(
      this.#gl.VERTEX_SHADER,
      viewportVertexShaderSource
    );
    const fragmentShader = this.#compileShader(
      this.#gl.FRAGMENT_SHADER,
      viewportFragmentShaderSource
    );

    this.#createProgram(vertexShader, fragmentShader);
    this.#createPosAttribute();
    this.#createUniforms();
    this.updateViewScale();
    this.#gl.clearColor(0.0, 0.0, 0.0, 0.0);
  }

  getCanvas() {
    return this.#cvs;
  }

  ok() {
    return this.#gl != null;
  }

  draw() {
    this.redrawCursor();
    this.#gl.bindFramebuffer(this.#gl.FRAMEBUFFER, null);
    this.#gl.clear(this.#gl.COLOR_BUFFER_BIT);
    this.#gl.drawArrays(this.#gl.TRIANGLES, 0, 6);
  }

  setTexture(img) {
    this.#gl.uniform2f(this.#u_size, img.width, img.height);
    this.#tex = this.#gl.createTexture();
    this.#gl.bindTexture(this.#gl.TEXTURE_2D, this.#tex);
    this.#gl.texParameteri(
      this.#gl.TEXTURE_2D,
      this.#gl.TEXTURE_WRAP_S,
      this.#gl.CLAMP_TO_EDGE
    );
    this.#gl.texParameteri(
      this.#gl.TEXTURE_2D,
      this.#gl.TEXTURE_WRAP_T,
      this.#gl.CLAMP_TO_EDGE
    );
    this.#gl.texParameteri(
      this.#gl.TEXTURE_2D,
      this.#gl.TEXTURE_MIN_FILTER,
      this.#gl.LINEAR
    );
    this.#gl.texParameteri(
      this.#gl.TEXTURE_2D,
      this.#gl.TEXTURE_MAG_FILTER,
      this.#gl.NEAREST
    );
    this.#gl.texImage2D(
      this.#gl.TEXTURE_2D,
      0,
      this.#gl.RGBA,
      this.#gl.RGBA,
      this.#gl.UNSIGNED_BYTE,
      img
    );
    this.#texFramebuffer = this.#gl.createFramebuffer();
    this.#gl.bindFramebuffer(this.#gl.FRAMEBUFFER, this.#texFramebuffer);
    this.#gl.framebufferTexture2D(
      this.#gl.FRAMEBUFFER,
      this.#gl.COLOR_ATTACHMENT0,
      this.#gl.TEXTURE_2D,
      this.#tex,
      0
    );
    this.#texScale = { x: img.width, y: img.height };
    this.#gl.uniform2f(this.#u_tex, this.#texScale.x, this.#texScale.y);
    if (this.#cvs.width > this.#cvs.height) {
      this.#zoom = this.#cvs.width / this.#texScale.x;
    } else {
      this.#zoom = this.#cvs.height / this.#texScale.y;
    }
    this.setZoom(this.#zoom);
  }

  setPixelColor(x, y, color) {
    let rgba = new Uint8Array(4);
    rgba[3] = 255;
    for (let i = 0; i < color.length; i++) {
      rgba[i] = color[i];
    }
    this.#gl.texSubImage2D(
      this.#gl.TEXTURE_2D,
      0,
      x,
      y,
      1,
      1,
      this.#gl.RGBA,
      this.#gl.UNSIGNED_BYTE,
      rgba
    );
  }

  setPixelBorder(x, y, color) {
    let rgba = new Uint8Array(4);
    rgba[3] = 254;
    for (let i = 0; i < color.length; i++) {
      rgba[i] = color[i];
    }
    // Draw the outline of the pixel using a hollow square
    this.#gl.texSubImage2D(
      this.#gl.TEXTURE_2D,
      0,
      x,
      y,
      1,
      1,
      this.#gl.RGBA,
      this.#gl.UNSIGNED_BYTE,
      rgba
    );
  }

  getColor(pos) {
    let rgba = new Uint8Array(4);
    this.#gl.bindFramebuffer(this.#gl.FRAMEBUFFER, this.#texFramebuffer);
    this.#gl.readPixels(
      pos.x,
      pos.y,
      1,
      1,
      this.#gl.RGBA,
      this.#gl.UNSIGNED_BYTE,
      rgba
    );
    return rgba.slice(0, 3);
  }

  getImageSize() {
    return this.#texScale;
  }

  scroll(ev) {
    this.#camPos = { x: ev.target.scrollLeft, y: ev.target.scrollTop };
    this.#gl.uniform2f(this.#u_cam, this.#camPos.x, this.#camPos.y);
  }

  move(x, y) {
    this.#camPos.x -= x / this.#zoom;
    this.#camPos.y -= y / this.#zoom;
    this.#gl.uniform2f(this.#u_cam, this.#camPos.x, this.#camPos.y);
  }

  getPos() {
    return this.#camPos;
  }

  setPos(x, y) {
    this.#camPos = { x: x, y: y };
    this.#gl.uniform2f(this.#u_cam, this.#camPos.x, this.#camPos.y);
  }

  transitionToPos(targetPos, duration, updateFunc) {
    const startTime = performance.now();
    const startPos = { x: this.#camPos.x, y: this.#camPos.y };
    let ctx = this;

    function easingFunc(t) {
      return t * t * (3.0 - 2.0 * t);
    }

    function animate(currentTime) {
      const elapsedTime = currentTime - startTime;
      const progress = Math.min(elapsedTime / duration, 1);
      const easedProgress = easingFunc(progress);
      const x = startPos.x + (targetPos.x - startPos.x) * easedProgress;
      const y = startPos.y + (targetPos.y - startPos.y) * easedProgress;
      ctx.setPos(x, y);
      ctx.draw();
      if (updateFunc) updateFunc(progress);
      if (progress < 1) {
        requestAnimationFrame(animate);
      }
    }

    requestAnimationFrame(animate);
  }

  setZoom(z) {
    if (z < 0.01) z = 0.01;
    if (z > 40) z = 40;
    this.#zoom = z;
    this.#gl.uniform1f(this.#u_zoom, z);
  }

  getZoom() {
    return this.#zoom;
  }

  updateViewScale() {
    let w = this.#cvs.clientWidth;
    let h = this.#cvs.clientHeight;
    this.#cvs.width = w;
    this.#cvs.height = h;
    this.#gl.viewport(0, 0, w, h);
    this.#gl.uniform2f(this.#u_view, w, h);
  }

  click(pos) {
    pos.x /= this.#cvs.width;
    pos.y /= this.#cvs.height;

    let a = {
      x:
        ((-0.5 * this.#texScale.x - this.#camPos.x) * this.#zoom) /
        this.#cvs.width +
        0.5,
      y:
        ((-0.5 * this.#texScale.y - this.#camPos.y) * this.#zoom) /
        this.#cvs.height +
        0.5,
    };

    let b = {
      x:
        ((0.5 * this.#texScale.x - this.#camPos.x) * this.#zoom) /
        this.#cvs.width +
        0.5,
      y:
        ((0.5 * this.#texScale.y - this.#camPos.y) * this.#zoom) /
        this.#cvs.height +
        0.5,
    };

    if (pos.x < a.x || pos.y < a.y || pos.x > b.x || pos.y > b.y) {
      return;
    }

    pos = {
      x: ((pos.x - a.x) / (b.x - a.x)) * this.#texScale.x,
      y: ((pos.y - a.y) / (b.y - a.y)) * this.#texScale.y,
    };

    return pos;
  }

  redrawCursor = () => {
    let color = this.getColor(this.outline);
    this.setPixelColor(this.outline.x + 0.5, this.outline.y + 0.5, color);
    try {
      let pos = this.click({
        x: window.innerWidth / 2,
        y: window.innerHeight / 2,
      });
      this.outline = { x: pos.x, y: pos.y };
    } catch (e) {
      return;
    }
    this.outline.x = Math.floor(this.outline.x);
    this.outline.y = Math.floor(this.outline.y);
    color = this.getColor(this.outline);
    this.setPixelBorder(this.outline.x, this.outline.y, color);
  };

  #createProgram(vertexShader, fragmentShader) {
    this.#program = this.#gl.createProgram();
    this.#gl.attachShader(this.#program, vertexShader);
    this.#gl.attachShader(this.#program, fragmentShader);
    this.#gl.linkProgram(this.#program);
    if (!this.#gl.getProgramParameter(this.#program, this.#gl.LINK_STATUS)) {
      console.error(this.#gl.getProgramInfoLog(this.#program));
      return null;
    }
    this.#gl.useProgram(this.#program);
  }

  #compileShader(type, source) {
    let shader = this.#gl.createShader(type);
    this.#gl.shaderSource(shader, source);
    this.#gl.compileShader(shader);
    if (!this.#gl.getShaderParameter(shader, this.#gl.COMPILE_STATUS)) {
      console.error(this.#gl.getShaderInfoLog(shader));
      this.#gl.deleteShader(shader);
      return null;
    }
    return shader;
  }

  #createPosAttribute() {
    let buffer = this.#gl.createBuffer();
    this.#gl.bindBuffer(this.#gl.ARRAY_BUFFER, buffer);
    let positions = [
      -0.5, -0.5, 0.5, -0.5, 0.5, 0.5, -0.5, -0.5, 0.5, 0.5, -0.5, 0.5,
    ];
    this.#gl.bufferData(
      this.#gl.ARRAY_BUFFER,
      new Float32Array(positions),
      this.#gl.STATIC_DRAW
    );
    this.#a_vert = this.#gl.getAttribLocation(this.#program, "vert");
    this.#gl.vertexAttribPointer(this.#a_vert, 2, this.#gl.FLOAT, false, 0, 0);
    this.#gl.enableVertexAttribArray(this.#a_vert);
  }

  #createUniforms() {
    this.#u_cam = this.#gl.getUniformLocation(this.#program, "cam");
    this.#u_tex = this.#gl.getUniformLocation(this.#program, "tex_scale");
    this.#u_view = this.#gl.getUniformLocation(this.#program, "view_scale");
    this.#u_zoom = this.#gl.getUniformLocation(this.#program, "zoom");
    this.#u_size = this.#gl.getUniformLocation(this.#program, "tex_size");
  }
}
