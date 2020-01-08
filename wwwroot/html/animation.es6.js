const animationAttribute = 'data-animation';
const defaultAvatarUrl = '/wwwroot/html/assets/logo.png';

function capitialize(str) {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

function getAvatarUrlByRtx(rtx) {
  return AVATAR[rtx];
}

class Animator {
  constructor($element) {
    this.$element = $element;
  }
}

class GoldPacketAnimator extends Animator {
  constructor($element) {
    super($element);
    this.opened = false;
  }

  close() {
    this.$element.removeClass('gold-packet-opened');
  }

  open() {
    this.$element.addClass('gold-packet-opened');
  }

  /***
   * 更新黄金红包的元素
   * {
   *    rtx: '发红包的人（默认Martinlau）',
   *    word: '祝福语（默认：恭贺新禧，日进斗金!）',
   *    total: '黄金红包总毫克',
   *    equalMoney: '等值金额',
   *    luckyStuffs: [{ avatar, rtx, packetAmount: '每个用户收到的黄金毫克数' }]
   * }
   * @param data
   */
  update(data) {
    const { $element, opened } = this;
    data.rtx = data.rtx || animator.goldPacket.rtx || 'martinlau';

    const capitalizedRtxName = capitialize(data.rtx);

    $element.find('.gold-packet-issuer-name').text(capitalizedRtxName);
    $element.find('.gold-packet-issuer').attr('src', getAvatarUrlByRtx(data.rtx) || defaultAvatarUrl);
      $element.find('.gold-packet-word').text(data.word || 'FiT AND BEST·2019');

    if (typeof data.total !== 'undefined') {
      $element.find('.gold-packet-amount-number').text(parseInt(data.total) || 0);
    }
    $element.find('.gold-packet-amount-measure-number').text(data.equalMoney ? data.equalMoney.toFixed(2) : undefined);

    const $detail = $element.find('.gold-packet-detail');
    $detail.html('');

    if (data.luckyStuffs && !opened) {
      data.luckyStuffs.forEach((stuff) => {
        const $node = $(`
            <div class="gold-packet-detail-row">
                <div>
                    <img class="gold-packet-detail-avatar" src=""/>
                </div>
                <span class="gold-packet-detail-name"></span>
                <span style="flex: 1"></span>
                <span class="gold-packet-detail-amount">
                    <span class="gold-packet-detail-amount-number"></span>
                    克
                </span>
            </div>
        `);

        $node.find('.gold-packet-detail-avatar').attr('src', stuff.avatar || defaultAvatarUrl);
        $node.find('.gold-packet-detail-name').text(stuff.rtx);
        $node.find('.gold-packet-detail-amount-number').text(stuff.packetAmount);

        $detail.append($node);
      });
    }
  }

  updateAndOpen(data) {
    return new Promise(resolve => {
      this.update(data);
      this.$element.one('transitionend', resolve);
      setTimeout(() => this.open());
    });
  }
}

class PageAnimator extends Animator {
  constructor($element) {
    super($element);
    this.goldPacketAnimator = new GoldPacketAnimator($('.gold-packet'));
    this.goldPacket =  false;
    this.packet = false;
    this.childAnimators = [];
    this.isRolling = false;
    this.childScale = 1;
    this.rollingPromiseList = [];
  }

  forEachAnimator(fn) {
    return this.childAnimators.filter(animator => animator.isShow).forEach(fn);
  }

  mapAnimator(fn) {
    return this.childAnimators.filter(animator => animator.isShow).map(fn);
  }

  /***
   * 获取下一次动画结束的Promise对象
   * @returns {Promise.<*[]>}
   */
  waitAnimationResolved(prevFn, selectors) {
    const { $element } = this;
    const promise = Promise.all(selectors.map((selector) => {
      return new Promise((resolve) =>
        $element.find(selector)
          .one('animationend', () => setTimeout(resolve)))
    }));
    setTimeout(prevFn);
    return promise;
  }


  /***
   * 更新DOM节点数据
   * @param data
   * {
     *   gift: {
     *     image: 'image_url',
     *     name: 'gift_name',
     *     desc: 'desc'
     *   },
     *   childRow: []
     * }
   */
  updateElements(data) {
    const { $element } = this;
    const $lotteryGift = $element.find('#lottery-gift');

    if (this.goldPacketAnimator) {
      this.goldPacketAnimator.close();
    }

    this.goldPacket = false;
    this.packet = false;

    if (data.goldPacket) {
      this.goldPacket = data.goldPacket;
      this.goldPacketAnimator.update(data.goldPacket);
      this.goldPacketAnimator.$element.show();
      $lotteryGift.hide();
    } else if (data.packet) {
      this.packet = data.packet;
      let title = `奖金池 ￥${(parseInt(data.packet.amount) || 0)}`;
      $lotteryGift.find('.lottery-gift-type').attr('src', data.packet.typeImg);
      $lotteryGift.find('.lottery-gift-image').attr('src', data.packet.image);
      $lotteryGift.find('.lottery-gift-name').html(title);
      $lotteryGift.find('.lottery-gift-desc').html('');
      this.goldPacketAnimator.$element.hide();
      $lotteryGift.show();
    } else {
      $lotteryGift.find('.lottery-gift-type').attr('src', data.gift.typeImg);
      $lotteryGift.find('.lottery-gift-image').attr('src', data.gift.image);
      $lotteryGift.find('.lottery-gift-name').html(data.gift.name);
      $lotteryGift.find('.lottery-gift-desc').html(data.gift.desc);
      this.goldPacketAnimator.$element.hide();
      $lotteryGift.show();
    }

    const $luckyStuffGroup = $element.find('.lucky-stuff-group');
    const $prevRows = $element.find('.lucky-stuff-row');

    $prevRows.children('.lucky-stuff-cloned').remove();
    $prevRows.each((i, element) => {
      $(element).addClass('hide')
    });

    this.childAnimators.forEach((animator) => {
      animator.isShow = false;
    });

    const $placeholders = [];
    let jj = 0;
    data.childRow.forEach((count, row) => {
      const rawPrevRow = $prevRows.get(row);
      const $row = rawPrevRow ? $(rawPrevRow) : $('<div class="lucky-stuff-row">');

      $row.children('.lucky-stuff').each((i, element) => {
        $(element).addClass('hide');
      });

      for (let i = 0; i < count; i++) {
        const cachedAnimator = this.childAnimators[jj];
        jj++;

        if (cachedAnimator) {
          cachedAnimator.updateElementFn(cachedAnimator.$element.find(`[${animationAttribute}=current]`), {
            avatar: defaultAvatarUrl,
            name: '虚席以待',
            rtx: ''
          });
          cachedAnimator.isShow = true;
          cachedAnimator.highlighted = false;
          cachedAnimator.returnMoneyAmount = 0;
          cachedAnimator.$element.removeClass('lucky-stuff-highlighted');
          cachedAnimator.$element.removeClass('hide');
          if (!$row.has(cachedAnimator.$element).length) {
            $row.append(cachedAnimator.$element);
          }
          continue;
        }

        const $placeholder = $(`
            <div class="lucky-stuff">
                <div class="lucky-stuff-avatar-container js-animation-scroll" data-store="avatar">
                    <div class="lucky-stuff-avatar js-animation-scroll-item">
                        <img src="${defaultAvatarUrl}"/>
                    </div>
                </div>

                <div class="lucky-stuff-information-container js-animation-scroll" data-store="text">
                    <div class="lucky-stuff-information js-animation-scroll-item">
                        <div class="lucky-stuff-name">虚席以待</div>
                        <div class="lucky-stuff-rtx"></div>
                    </div>
                </div>
            </div>
            `);

        $row.append($placeholder);
        $placeholders.push($placeholder);
      }

      if (!rawPrevRow) {
        $luckyStuffGroup.append($row);
      } else {
        $row.removeClass('hide')
      }
    });

    $placeholders.forEach(($placeholder) => {
      this.childAnimators.push(
        new LuckyStuffAnimator($placeholder, ($node, _data) => {
          $node.find('img').attr('src', getAvatarUrlByRtx(_data.rtx) || defaultAvatarUrl);
          $node.find('.lucky-stuff-name').text(_data.name);
          $node.find('.lucky-stuff-rtx').text(_data.rtx);
        }, ['text', 'avatar'])
      );
    });

    return new Promise((resolve) => {
      const $groupContainer = $luckyStuffGroup.closest('.lucky-stuff-group-container');
        $luckyStuffGroup.parent('.lucky-stuff-group-js-placeholder')
            .css('height', $luckyStuffGroup.height());

        //const widthScale = $groupContainer.width() / $luckyStuffGroup.width();
        const widthScale = 1300 / $luckyStuffGroup.width();
        const heightScale  = $groupContainer.height() / $luckyStuffGroup.height();

        const scale = this.childScale = Math.min(widthScale, heightScale, 1.75);
        $groupContainer.css({ transform: `scale(${scale})`});

        setTimeout(resolve);
    });
  }

  /***
   * 设置当前节点状态以触发动画
   * @param state
   *  show 显示奖品并置于中间 // 黄金红包使用
   *  before-draw 奖品左移 并显示抽奖列表 // 黄金红包使用
   *  before-draw-direct 直接显示奖品与抽奖列表
   *  hide 隐藏所有节点 // 用于更新DOM
   *  before-hide 执行隐藏之前的动画
   *  hide-center 隐藏居中 // 黄金红包
   */
  setCurrentState(state) {
    this.$element.attr('data-state', state);
  }
}

class LuckyStuffAnimator extends Animator {
  constructor($element, updateElementFn, animateTypes) {
    super($element);
    this.updateElementFn = updateElementFn;
    this.animateTypes = animateTypes;
    this.returnMoneyAmount = 0;
    this.isShow = true;

    setTimeout(() => this.initialize());
  }

  initialize() {
    const {$element} = this;
    const $scroll = $element.find('.js-animation-scroll');

    $scroll.each((index, element) => {
      const $scrollElement = $(element);

      const width = $scrollElement.width();
      const height = $scrollElement.height();

      /* set height for container */
      $scrollElement.width(width);
      $scrollElement.height(height);

      const $node1 = $scrollElement.find('.js-animation-scroll-item');
      const $node2 = $node1.clone();
      const $node3 = $node1.clone();

      /* add data attribute to control style */
      $node1.attr(animationAttribute, 'current');
      $node2.attr(animationAttribute, 'prev');
      $node3.attr(animationAttribute, 'next');

      $scrollElement.append($node2);
      $scrollElement.append($node3);

      this[$scrollElement.attr('data-store')] = {
        $prev: $node2,
        $next: $node3,
        $current: $node1
      };
    });
  }

  animate() {
    const promise = new Promise((resolve) => {
      this.$element.one('transitionend', () => setTimeout(resolve));
    });

    this.animateTypes.forEach((type) => {
      const {$prev, $next, $current} = this[type];

      this[type] = {
        $prev: $current.attr(animationAttribute, 'prev'),
        $next: $prev.attr(animationAttribute, 'next'),
        $current: $next.attr(animationAttribute, 'current')
      };
    });

    return promise;
  };

  updateElement(data) {
    const { $element } = this;
    const $nodes = $element.find(`[${animationAttribute}=next]`);
    if (data.isLeader) {
      this.highlighted = true;
      this.returnMoneyAmount = data.backmoney;
      $element.addClass('lucky-stuff-highlighted');
    } else {
      this.highlighted = false;
      this.returnMoneyAmount = 0;
      $element.removeClass('lucky-stuff-highlighted');
    }
    this.updateElementFn($nodes, data);
  }

  /***
   * 根据传入数据更新DOM并滚动节点文字
   * @param data 格式 { image, name, rtx }
   * @returns {Promise}
   */
  updateAndAnimate(data) {
    return new Promise((resolve) => {
      this.updateElement(data);
      requestAnimationFrame(() => {
        this.animate().then(resolve);
      });
    });
  }
}

class LuckyBagAnimator extends Animator {
  constructor($element, currentNumber) {
    super($element);
    this.currentNumber = currentNumber || 0;
    this.initialize();
  }

  initialize() {
    const { $element } = this;
    this.$counter = $element.find('.lucky-bag-counter');
    this.$increase = $element.find('.lucky-bag-increase');

    this.$counter.addClass('odometer');
  }

  hide() {
    const { $element } = this;
    $element.hide();
  }

  show() {
    const { $element, $increase } = this;
    $increase.removeClass('animating');
    $element.show();
  }

  update(value) {
    const { $counter } = this;
    if (value === this.currentNumber) return false;
    const increment = value - this.currentNumber;
    let isIncrease = increment >= 0;
    this.currentNumber = value;
    $counter.html(value);
    this.showIncrement(Math.abs(increment), !isIncrease);
  }

  showIncrement(increment, reverse) {
    const { $increase } = this;
    $increase.removeClass('animating');
    if (reverse) {
      $increase.removeClass('increase');
      $increase.addClass('decrease');
    } else {
      $increase.addClass('increase');
      $increase.removeClass('decrease');
    }
    setTimeout(() => {
      $increase.html(increment);
      $increase.addClass('animating');
    });
  }
}

const animator = new PageAnimator($('#container'));
const luckyBagAnimator = new LuckyBagAnimator($('.lucky-bag'));

/***
 *
 * @param data
 * 普通抽奖 { gift: {} }
 * 黄金红包 { goldPacket: { rtx, word, total?, equalMoney? } }
 * @returns {Promise.<*[]>}
 */
function ready(data) {
  let promise = hide()
    .then(() => setTimeout(() => animator.updateElements(data)));

  if (data.goldPacket) {
    promise = promise.then(() => {
      return animator.waitAnimationResolved(() => {
        return animator.setCurrentState('show');
      }, ['.gold-packet']).then(() => {
        return new Promise((resolve) => {
          /* 黄金红包 停留seconds */
          setTimeout(resolve, 1000);
        }).then(() => {
          return animator.waitAnimationResolved(() => {
            return animator.setCurrentState('before-draw');
          }, ['.gold-packet', 'lucky-stuff']);
        });
      });
    });
  } else {
    promise = promise.then(() => {
      return animator.waitAnimationResolved(() => {
        return animator.setCurrentState('before-draw-direct');
      }, ['#lottery-gift', 'lucky-stuff']);
    });
  }
  animator.isRolling = false;
  return promise;
}

/***
 * 隐藏当前抽奖
 * @returns {*}
 */
function hide() {
  const { $element } = animator;
  if ($element.attr('data-state')) {
    const state = $element.attr('data-state');
    const selectors = [];
    switch (selectors) {
      case 'before-draw':
      case 'before-draw-direct':
        selectors.push('.lucky-stuff-group');
        break;
    }

    let hideState = state === 'before-hide' ? 'hide-center' : 'hide';
    if (animator.goldPacket) {
      selectors.push('.gold-packet');
    } else {
      selectors.push('#lottery-gift');
    }

    return animator
      .waitAnimationResolved(() => animator.setCurrentState(hideState), selectors)
      .then(() => setTimeout(() => animator.setCurrentState('')));
  } else {
    return Promise.resolve();
  }
}

/***
 * @params userCollection [{avatar: string, name: string, rtx: string}]
 * 开始抽奖
 */
function start(userCollection) {
  if (animator.packet && (animator.packet.rtx.toLowerCase() === 'pool' || animator.packet.rtx === '现金池')) {
    const cost = parseInt(animator.packet.amount) * parseInt(animator.packet.count);
    const updatedValue = luckyBagAnimator.currentNumber - cost;
    luckyBagAnimator.update(updatedValue);
  }
  animator.isRolling = true;
  if (userCollection.length) {
    return roll(userCollection);
  } else {
    return Promise.reject();
  }
}

/***
 * 当前滚动对象
 * @params userCollection [{avatar: string, name: string, rtx: string}]
 * @returns {Promise.<*[]>}
 */
function roll(userCollection) {
  const promises = animator.mapAnimator((childAnimator) => {
    return rollSingleItem(childAnimator, userCollection).then();
  });
  let done = Promise.all(promises);
  animator.rollingPromiseList = promises;
  return done;
}

function rollSingleItem(childAnimator, userCollection) {
  const len = userCollection.length;
  return childAnimator.updateAndAnimate(userCollection[parseInt(Math.random() * (len - 1))]).then(() => {
    if (animator.isRolling) {
      return rollSingleItem(childAnimator, userCollection);
    } else {
      return Promise.resolve();
    }
  });
}

/***
 * 显示最终抽奖结果并结束动画
 * @param data [{avatar: string, name: string, rtx: string}]
 * @returns {Promise.<>}
 */
function showResult(data) {
  const promises = animator.mapAnimator((childAnimator, i) => {
    const roll = animator.rollingPromiseList[i] || Promise.resolve();
    roll.then(() => childAnimator.updateAndAnimate(data[i]));
  });

  let promise = Promise.all(promises);

  promise = promise.then(() => {
    return new Promise((resolve) => {
      setTimeout(resolve, 1500);
    });
  });

  if (animator.goldPacket) {
    const goldPacketData = {
      luckyStuffs: []
    };

    const packetAmount = (animator.goldPacket.total / data.length || 0).toFixed(2);
    data.forEach((stuff) => {
      goldPacketData.luckyStuffs.push({
        rtx: stuff.rtx,
        avatar: getAvatarUrlByRtx(stuff.rtx),
        packetAmount: packetAmount
      })
    });

    animator.goldPacketData = goldPacketData;
  }

  animator.isRolling = false;
  return promise;
}

/***
 * 把Leader放入钱袋
 * @returns {Promise.<>}
 */
function hideHighlighted() {
  const { offsetWidth, offsetHeight } = document.body;
  const luckyBagRect = $('.lucky-bag').get(0).getBoundingClientRect();
  const currentScale = animator.childScale;

  const center = {
    x: offsetWidth / 2,
    y: offsetHeight / 2
  };

  const bagCenter = {
    x: luckyBagRect.x + luckyBagRect.width / 2,
    y: luckyBagRect.y + luckyBagRect.height / 2
  };

  let promise = Promise.resolve();

  animator.mapAnimator((animator) => {
    if (!animator.highlighted || !animator.returnMoneyAmount) return;
    const { $element } = animator;
    const rawElement = $element.get(0);
    const rect = rawElement.getBoundingClientRect();
    const $clonedElement = $element.clone();

    $clonedElement.addClass('lucky-stuff-cloned');

    $clonedElement.css({
      position: 'absolute',
      left: rawElement.offsetLeft - 18,
      /* 减去margin-top */
      top: rawElement.offsetTop - 18
    });

    $clonedElement.width($element.width());
    $clonedElement.find('.lucky-stuff-information-container').hide();

    $clonedElement.insertAfter($element);

    const textHeight = $clonedElement
      .find('.lucky-stuff-information-container')
      .height();

    const rectCenter = {
      x: rect.x + rect.width / 2,
      y: rect.y + rect.height / 2 - textHeight
    };

    const centerOffsetX = center.x - rectCenter.x;
    const centerOffsetY = center.y - rectCenter.y;

    promise = promise
      .then(() => {
        const promise = new Promise((resolve) => {
          $clonedElement.one('transitionend', (event) => {
            setTimeout(resolve, 500);
          });
        });

        const actualWidth = $clonedElement.get(0).getBoundingClientRect().width;

        $clonedElement.addClass('lucky-stuff-zooming');
        $clonedElement.css({
          transform: `
          translate3d(${centerOffsetX / currentScale}px,
           ${centerOffsetY / currentScale}px, 0) scale(${450 / actualWidth})`,
          transitionDuration: '0.5s',
          zIndex: 10
        });

        return promise;
      })
      .then(() => {
        const bagOffsetX = bagCenter.x - rectCenter.x;
        const bagOffsetY = bagCenter.y - rectCenter.y;

        const promise = new Promise((resolve) => {
          $clonedElement.one('transitionend', () => setTimeout(resolve));
        });

        $clonedElement.css({
          transform: `
          translate3d(${bagOffsetX / currentScale}px,
           ${bagOffsetY / currentScale}px, 0)`,
          transitionDuration: '0.2s',
          opacity: 0
        });

        return promise;
      }).then(() => luckyBagAnimator.update(luckyBagAnimator.currentNumber + animator.returnMoneyAmount));
  });

  return promise;
}

function beforeHide() {
  let promise = hideHighlighted();

  if (animator.goldPacket) {
    promise = promise.then(() => {
      return animator.waitAnimationResolved(() => {
        animator.setCurrentState('before-hide');
      }, ['.gold-packet', '.lucky-stuff-group'])
        .then(() => {
          animator.goldPacketAnimator.updateAndOpen(animator.goldPacketData)
        });
    });
  }

  return promise;
}

/***
 * 更新钱袋数字
 * @param value
 */
function updateBag(value) {
  luckyBagAnimator.update(parseInt(value, 10));
}

/***
 * 隐藏奖池数字
 * @param value
 */
function hideBag() {
  luckyBagAnimator.hide();
}

/***
 * 显示奖池数字
 * @param value
 */
function showBag() {
  luckyBagAnimator.show();
}

function toggleMainView(show) {
  return hide().then(() => {
    if (show) {
      $('.mask-view').removeClass('mask-view-hide');
    } else {
      $('.mask-view').addClass('mask-view-hide');
    }
  });
}
