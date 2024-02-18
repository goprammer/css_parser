package css_parser

const TestStr = `body{margin:0;padding:0;}img { -ms-interpolation-mode:bicubic; border:0; outline:none; max-width: 100%;margin: 10px 13px 15px !important;}
	#img{margin:13px}.img{margin-left:3px/*Another Comment*/;margin-right:1px;}
	#comment /*Comment*/
				{
					margin: 21px
				}
	p#complex {color:dodgerblue;}
	#complex {color:red}
	.inlineBlock {display: inline-block !important;} 
	.block {display: block;} .tableCell {display: table-cell !important;}
	.align_left {text-align: left !important;} 
	.align_right {text-align: right !important;;}
	.w40 {width:40%} 
	.w49 {width:49%} 
	.w50 {width:50%} 
	.w59 {width:59%} 
	.mw_100 {max-width: 100%;} 
	.w_290 {width: 290px;}
	.h_200 {height: 200px !important; }
	.floatRight { float: right !important; }
	.desktop_hide, #laptop_hide {mso-hide: all; display: none; max-height: 0px; overflow: hidden;}
	.plr0 {padding-left: 0px !important; padding-right: 0px !important;} 
	.plr10 {padding-left: 10px !important; padding-right: 10px !important;}
	@media (max-width: 620px) {
		body {padding-left: 10px !important;padding-right: 10px !important;}
		.mobile_full {width: 100% !important;min-width: 100% !important;padding: 0px 0px 0px 0px !important;}
		.mobile_full100 {width: 100% !important;min-width: 100% !important;}
		.mobile_hide {min-height: 0;max-height: 0;max-width: 0;display: none !important;overflow: hidden;font-size: 0;}
		.desktop_hide {display: table-cell !important;max-height: none !important;}
		.mobile_inline {display: inline !important;}
		.mobile_block {display: block !important;}
		.mobile_inlineBlock {display: inline-block !important;}
		.mobile_maxMinAuto {max-width: auto !important;min-width: auto !important;max-height: auto !important;min-height: auto !important;}
		.mobile_minAuto {min-width: auto !important;}
		.mobile_widthAuto {width: auto !important;}
		.img-container.big img {width: auto !important;}
		.mobile_90 {width: 90% !important;}
		.mobile_10 {width: 10% !important;}
		.mobile_80 {width: 80% !important;}
		.mobile_20 {width: 20% !important;}
		.mobile_84 {width: 84% !important;}
		.mobile_16 {width: 16% !important;}
		.mobile_30 {width: 30% !important;}
		.mobile_38 {width: 38% !important;}
		.mobile_42 {width: 42% !important;}
		.mobile_46 {width: 46% !important;}
		.mobile_49 {width: 49% !important;}
		.mobile_50 {width: 50% !important;}
		.mobile_53 {width: 53% !important;}
		.mobile_57 {width: 57% !important;}
		.mobile_60 {width: 60% !important;}
		.mobile_70 {width: 70% !important;}
		.w100 {width: 100% !important;}
		mobile_w18 {width: 18px !important;}
		.mobile_w41 { /* etv */width: 41px !important;}
		.mobile_alignleft {text-align: left !important;}
		.mobile_alignright {text-align: right !important;}
		.mobile_aligncenter {text-align: center !important;}
		.mobile_floatRight {float: right !important;}
		.mobile_mauto {margin: auto !important;}
		.mobile_plr5 {padding-left: 5px !important;padding-right: 5px !important;}
		.mobile_plr10 {padding-left: 10px !important;padding-right: 10px !important;}
		.mobile_plr16 {padding-left: 16px !important;padding-right: 16px !important;}
		.mobile_ptb20 {padding-top: 20px !important;padding-bottom: 20px !important;}
		.mobile_ptb28 {padding-top: 28px !important;padding-bottom: 28px !important;}
		.mobile_pt10 {padding-top: 10px !important;}
		.mobile_pt12 {padding-top: 12px !important;}
		.mobile_pt16 {padding-top: 16px !important;}
		.mobile_pt18 {padding-top: 18px !important;}
		.mobile_plr0 {padding-left: 0px !important;padding-right: 0px !important;}
		.mobile_pt12b17 {padding-top: 12px !important;padding-bottom: 17px !important;}
		.mobile_px37y15 {padding-left: 37px !important;padding-right: 37px !important;padding-top: 15px !important;padding-bottom: 15px !important;}
		.mobile_px30y11 { /* etv */padding-left: 30px !important;padding-right: 30px !important;padding-top: 11px !important;padding-bottom: 11px !important;}
		.mobile_prefixPadding {padding: 2px 6px 0px 6px !important;}
		.mobile_pt0 {padding-top: 0px !important;}
		.mobile_pt8 {padding-top: 8px !important;}
		.mobile_pb0 {padding-bottom: 0px !important;}
		.mobile_p0 {padding: 0px !important;}
		.mobile_pl0 {padding-left: 0px !important;}
		.mobile_pr0 {padding-right: 0px !important;}
		.mobile_pr5 {padding-right: 5px !important;}
		.mobile_p10 {padding: 10px !important;}
		.mobile_pl15 {padding-left: 15px !important;}
		.mobile_pl20 {padding-left: 20px !important;}
		.mobile_pr10 {padding-right: 10px !important;}
		.mobile_pr20 {padding-right: 20px !important;}
		.mobile_plr25 {padding-left: 25px !important;padding-right: 25px !important;}
		.mobile_pb8 {padding-bottom: 8px !important;}
		.mobile_pb10 {padding-bottom: 10px !important;}
		.mobile_pb12 {padding-bottom: 12px !important;}
		.mobile_pb30 {padding-bottom: 30px !important;}
		.mobile_pb16 {padding-bottom: 16px !important;}
		.mobile_ptl15 {padding-top: 15px !important;padding-left: 15px !important;}
		.mobile_pt30 {padding-top: 30px !important;}
		.mobile_pt20 {padding-top: 20px !important;padding-bottom: 0px !important;}
		.mobile_pt15r15b25l15 {padding-top: 15px !important;padding-right: 15px !important;padding-bottom: 25px !important;padding-left: 15px !important;}
		.mb_fz12, .mobile-fs12 {font-size: 12px !important;line-height: 14px !important;}
		.mb_fz13, .mobile-fs13 {font-size: 13px !important;line-height: 15px !important;}
		.mb_fz15, .mobile-fs15 {font-size: 15px !important;line-height: 18px !important;}
		.mb_fz17, .mobile-fs17 {font-size: 17px !important;line-height: 19px !important;}
		.mb_fz18 {font-size: 18px !important;}
		.mb_lh21 {line-height: 21px !important;}
		.mb_fz20, .mobile-fs20 {font-size: 20px !important;line-height: 23px !important;}
		.mb_fz22 {font-size: 22px !important;}.mb_lh25 {line-height: 25px !important;}
		.mb_fz24, .mobile-fs24 {font-size: 24px !important;line-height: 27px !important;}
		.mb_fz27, .mobile-fs27 {font-size: 27px !important;line-height: 30px !important;}
		.authorRight {min-width: 63% !important;max-width: 63% !important;}
		.footer-share .td {display: block !important;width: 100% !important;}
		.footer-share .tdfirst {padding: 0 0 20px 0 !important;}
		.footer-share .tdmiddle {border-left: none !important;border-top: 1px solid #DBDBDB;padding: 20px 0 20px 0 !important;}
		.footer-share .tdlast {border-left: none !important;border-top: 1px solid #DBDBDB;padding: 20px 0 0 0 !important;}}
		@media (max-width:800px){
						.img{margin-left:11px
					}
				}

				@media(max-width:400px){.img{margin-left:18px}}
				@media (min-width:510px) and (max-width:512px) {
					.double {
						position:fixed;
					}
				}`