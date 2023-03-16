export const getMeta = (url) => {
  // Return width and height of image
  return new Promise((resolve, reject) => {
    var img = new Image();
    img.addEventListener("load", function () {
      let w = this.naturalWidth;
      let h = this.naturalHeight;
      resolve({ w, h });
    });
    img.addEventListener("error", function () {
      reject("Error loading image");
    });
    img.src = url;
  });
};
