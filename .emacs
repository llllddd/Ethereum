;; from https://www.youtube.com/watch?v=r6j2W5DZRtA
;; get the following packages ("M-x package-list-packages"):
;;     go-mode
;;     go-eldoc
;;     company-mode
;;     company-go
;; get the following go programs (run each line in your shell):
;;     go get golang.org/x/tools/cmd/godoc
;;     go get golang.org/x/tools/cmd/goimports
;;     go get github.com/rogpeppe/godef
;;     go get github.com/nsf/gocode

;;Cancel auto save
;;(setq auto-save-default nil)

;;Cancel auto backup
;;(setq auto-backup-files nil)

(require 'package)
(add-to-list 'package-archives
         '("melpa" . "http://melpa.milkbox.net/packages/") t)

(setq company-idle-delay t)

(setq gofmt-command "goimports")
;; UPDATE: gofmt-before-save is more convenient then having a command
;; for running gofmt manually. In practice, you want to
;; gofmt/goimports every time you save anyways.

(add-hook 'before-save-hook 'gofmt-before-save)

(global-set-key (kbd "C-c M-n") 'company-complete)
(global-set-key (kbd "C-c C-n") 'company-complete)

(defun my-go-mode-hook ()
  ;; UPDATE: I commented the next line out because it isn't needed
  ;; with the gofmt-before-save hook above.
  (local-set-key (kbd "C-c m") 'gofmt)
  (local-set-key (kbd "M-.") 'godef-jump)
  (set (make-local-variable 'company-backends) '(company-go)))

(add-hook 'go-mode-hook 'my-go-mode-hook)
(add-hook 'go-mode-hook 'go-eldoc-setup)
(add-hook 'go-mode-hook 'company-mode)


 ;;(custom-set-variables
 ;; custom-set-variables was added by Custom.
 ;; If you edit it by hand, you could mess it up, so be careful.
 ;; Your init file should contain only one such instance.
 ;; If there is more than one, they won't work right.
 ;; '(package-selected-packages '(go-autocomplete go-eldoc company-go)))
 ;;(custom-set-faces
 ;; custom-set-faces was added by Custom.
 ;; If you edit it by hand, you could mess it up, so be careful.
 ;; Your init file should contain only one such instance.
 ;; If there is more than one, they won't work right.
 ;;)

(global-linum-mode 1) ; always show line numbers                              
(setq linum-format "%d| ")  ;set format

;; flymake mode is used to check the synax error 
(add-hook 'flymake-mode-hook
      (lambda()
        (local-set-key (kbd "C-c C-e n") 'flymake-goto-next-error)))
(add-hook 'flymake-mode-hook
      (lambda()
        (local-set-key (kbd "C-c C-e p") 'flymake-goto-prev-error)))
(add-hook 'flymake-mode-hook
      (lambda()
        (local-set-key (kbd "C-c C-e m") 'flymake-popup-current-error-menu)))

;; change the theme
(load-theme 'monokai t)

(company-quickhelp-mode)
;;(setq company-quickhelp-delay t)
(eval-after-load 'company
  '(define-key company-active-map (kbd "C-c h") #'company-quickhelp-manual-begin))

;;(add-to-list 'load-path "~/.emacs.d/")
;;(require 'go-autocomplete)
;;(require 'auto-complete-config)
;;(ac-config-default)
;;Also you could setup any combination (for example M-TAB)
;;for invoking auto-complete:
